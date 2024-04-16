use std::collections::HashMap;
use std::net::{TcpListener, TcpStream};
use std::thread;
use std::io::{Write, BufReader, BufRead};
use std::sync::{Arc, Mutex};
use log::error;

enum Command {
    CapLs302,
    Nick(String),
    User(String),
    Quit,
    Join(String),
    Part(String),
    PrivMsg(String, String),
    Op(String),
    Deop(String),
    Ping,
    Pong,
    Names(String),
    Unknown(String),
}

fn parse_command(line: &str) -> Command {
    let parts: Vec<&str> = line.split_whitespace().collect();
    match parts.get(0) {
        Some(&"CAP") if parts.get(1) == Some(&"LS") && parts.get(2) == Some(&"302") => Command::CapLs302,
        Some(&"NICK") => Command::Nick(parts.get(1).cloned().unwrap_or("unknown").to_string()),
        Some(&"USER") => Command::User(parts.get(1).cloned().unwrap_or("unknown").to_string()),
        Some(&"QUIT") => Command::Quit,
        Some(&"JOIN") => Command::Join(parts.get(1).cloned().unwrap_or("unknown").to_string()),
        Some(&"PART") => Command::Part(parts.get(1).cloned().unwrap_or("unknown").to_string()),
        Some(&"PRIVMSG") if parts.len() >= 3 => {
            let target = parts[1].to_string();
            let message = parts[2..].join(" ");
            Command::PrivMsg(target, message)
        }
        Some(&"OP") => Command::Op(parts.get(1).cloned().unwrap_or("unknown").to_string()),
        Some(&"DEOP") => Command::Deop(parts.get(1).cloned().unwrap_or("unknown").to_string()),
        Some(&"PING") => Command::Ping,
        Some(&"PONG") => Command::Pong,
        Some(&"NAMES") => Command::Names(parts.get(1).cloned().unwrap_or("unknown").to_string()),
        _ => Command::Unknown(line.to_string()),
    }
}

fn handle_client(mut stream: TcpStream, channels: Arc<Mutex<HashMap<String, Vec<String>>>>, ops: Arc<Mutex<HashMap<String, String>>>, clients: Arc<Mutex<HashMap<String, TcpStream>>>) -> std::io::Result<()> {
    let welcome_msg = "Welcome to the IRC server!\r\n";
    stream.write_all(welcome_msg.as_bytes())?;

    let mut reader = BufReader::new(stream.try_clone()?);
    let mut buffer = String::new();
    let mut nickname = String::new();

    while reader.read_line(&mut buffer)? > 0 {
        let line = buffer.trim();
        let command = parse_command(line);
        match command {
            Command::CapLs302 => {
                stream.write_all("CAP * LS :\r\n".as_bytes())?;
            }
            Command::Nick(nick) => {
                nickname = nick.clone();
                clients.lock().unwrap().insert(nickname.clone(), stream.try_clone()?);
                stream.write_all(format!("NICK {}\r\n", nick).as_bytes())?;
            }
            Command::User(username) => {
                stream.write_all(format!("USER {}\r\n", username).as_bytes())?;
            }
            Command::Quit => {
                stream.write_all("Goodbye!\r\n".as_bytes())?;
                clients.lock().unwrap().remove(&nickname);
                return Ok(());
            }
            Command::Join(channel) => {
                let mut channels = channels.lock().unwrap();
                channels.entry(channel.clone()).or_insert_with(Vec::new).push(nickname.clone());
                let join_msg = format!(":{} JOIN {}\r\n", nickname, channel);
                for (_, client_stream) in clients.lock().unwrap().iter_mut() {
                    client_stream.write_all(join_msg.as_bytes())?;
                }
            }
            Command::Part(channel) => {
                let mut channels = channels.lock().unwrap();
                if let Some(members) = channels.get_mut(&channel) {
                    members.retain(|nick| nick != &nickname);
                    let part_msg = format!(":{} PART {}\r\n", nickname, channel);
                    for (_, client_stream) in clients.lock().unwrap().iter_mut() {
                        client_stream.write_all(part_msg.as_bytes())?;
                    }
                } else {
                    stream.write_all(format!("You're not in channel {}\r\n", channel).as_bytes())?;
                }
            }
            Command::PrivMsg(target, message) => {
                if let Some(target_stream) = clients.lock().unwrap().get(&target) {
                    target_stream.write_all(format!(":{} PRIVMSG {} :{}\r\n", nickname, target, message).as_bytes())?;
                } else {
                    stream.write_all(format!("No such nick/channel: {}\r\n", target).as_bytes())?;
                }
            }
            Command::Op(channel) => {
                let mut ops = ops.lock().unwrap();
                ops.insert(channel.clone(), nickname.clone());
                stream.write_all(format!(":{} OP {}\r\n", nickname, channel).as_bytes())?;
            }
            Command::Deop(channel) => {
                let mut ops = ops.lock().unwrap();
                if let Some(op) = ops.get(&channel) {
                    if *op == nickname {
                        ops.remove(&channel);
                        stream.write_all(format!(":{} DEOP {}\r\n", nickname, channel).as_bytes())?;
                    } else {
                        stream.write_all(format!("You're not an operator of channel {}\r\n", channel).as_bytes())?;
                    }
                } else {
                    stream.write_all(format!("No operator for channel {}\r\n", channel).as_bytes())?;
                }
            }
            Command::Ping => {
                stream.write_all("PONG\r\n".as_bytes())?;
            }
            Command::Pong => {
                // Do nothing
            }
            Command::Names(channel) => {
                let channels = channels.lock().unwrap();
                if let Some(members) = channels.get(&channel) {
                    let names = members.join(" ");
                    stream.write_all(format!(":{} NAMES {} :{}\r\n", nickname, channel, names).as_bytes())?;
                } else {
                    stream.write_all(format!("No such channel: {}\r\n", channel).as_bytes())?;
                }
            }
            Command::Unknown(cmd) => {
                stream.write_all(format!("Unknown command: {}\r\n", cmd).as_bytes())?;
            }
        }
        buffer.clear();
    }

    Ok(())
}

fn main() -> std::io::Result<()> {
    env_logger::init();
    let listener = TcpListener::bind("127.0.0.1:6667")?;
    let channels: Arc<Mutex<HashMap<String, Vec<String>>>> = Arc::new(Mutex::new(HashMap::new()));
    let ops: Arc<Mutex<HashMap<String, String>>> = Arc::new(Mutex::new(HashMap::new()));
    let clients: Arc<Mutex<HashMap<String, TcpStream>>> = Arc::new(Mutex::new(HashMap::new()));
    for stream in listener.incoming() {
        match stream {
            Ok(stream) => {
                let channels = Arc::clone(&channels);
                let ops = Arc::clone(&ops);
                let clients = Arc::clone(&clients);
                thread::spawn(move || {
                    if let Err(e) = handle_client(stream, channels, ops, clients) {
                        error!("Error handling client: {}", e);
                    }
                });
            }
            Err(e) => {
                error!("Failed: {}", e);
            }
        }
    }
    Ok(())
}