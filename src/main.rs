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
    Mode(String, String, String), // Added Mode command
    Who(String), // Added Who command
    Unknown(String),
}

fn parse_command(line: &str) -> Command {
    let parts: Vec<&str> = line.split_whitespace().collect();
    match parts.get(0) {
        Some(command) => match command.to_uppercase().as_str() {
            "CAP" if parts.get(1).map(|s| s.eq_ignore_ascii_case("LS")) == Some(true) && parts.get(2).map(|s| s.eq_ignore_ascii_case("302")) == Some(true) => Command::CapLs302,
            "NICK" => Command::Nick(parts.get(1).cloned().unwrap_or("unknown").to_string()),
            "USER" => Command::User(parts.get(1).cloned().unwrap_or("unknown").to_string()),
            "QUIT" => Command::Quit,
            "JOIN" => Command::Join(parts.get(1).cloned().unwrap_or("unknown").to_string()),
            "PART" => Command::Part(parts.get(1).cloned().unwrap_or("unknown").to_string()),
            "PRIVMSG" if parts.len() >= 3 => {
                let target = parts[1].to_string();
                let message = parts[2..].join(" ");
                Command::PrivMsg(target, message)
            }
            "MODE" if parts.len() >= 3 => {
                let channel = parts[1].to_string();
                let mode = parts[2].to_string();
                let user = parts.get(3).cloned().unwrap_or("unknown").to_string();
                Command::Mode(channel, mode, user)
            }
            "WHO" if parts.len() >= 2 => {
                let channel = parts[1].to_string();
                Command::Who(channel)
            }
            "OP" => Command::Op(parts.get(1).cloned().unwrap_or("unknown").to_string()),
            "DEOP" => Command::Deop(parts.get(1).cloned().unwrap_or("unknown").to_string()),
            "PING" => Command::Ping,
            "PONG" => Command::Pong,
            "NAMES" => Command::Names(parts.get(1).cloned().unwrap_or("unknown").to_string()),
            _ => Command::Unknown(line.to_string()),
        },
        None => Command::Unknown(line.to_string()),
    }
}

fn handle_client(mut stream: TcpStream, channels: Arc<Mutex<HashMap<String, Vec<String>>>>, ops: Arc<Mutex<HashMap<String, String>>>) -> std::io::Result<()> {
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
            Command::Ping => {
                stream.write_all("PONG\r\n".as_bytes())?;
            }
            Command::Pong => {
                // Do nothing
            }
            Command::Nick(nick) => {
                nickname = nick.clone();
                let nicknames: Arc<Mutex<HashSet<String>>> = Arc::new(Mutex::new(HashSet::new()));
                let hostname = "localhost"; // replace with your hostname
                let version = "0.1.0"; // replace with your version
                let date = "2022-01-01"; // replace with your server's creation date
                let servername = "myserver"; // replace with your server name
                let user_modes = "iws"; // replace with your available user modes
                let channel_modes = "imn"; // replace with your available channel modes
                let server_name = "localhost"; // replace with your server name
                let port_number = "6667"; // replace with your port number
                let motd_text = "Welcome to my IRC server!"; // replace with your message of the day
                thread::spawn(move || {
                    if let Err(e) = handle_client(stream, channels, ops, nicknames) {
                        error!("Error handling client: {}", e);
                    }
                });
                stream.write_all(format!("NICK {}\r\n", nick).as_bytes())?;
                stream.write_all(format!(":{} 001 {} :Welcome to the network, {}\r\n", nickname, nickname, nickname).as_bytes())?;
                stream.write_all(format!(":{} 002 {} :Your host is {}, running version {}\r\n", nickname, nickname, hostname, version).as_bytes())?;
                stream.write_all(format!(":{} 003 {} :This server was created {}\r\n", nickname, nickname, date).as_bytes())?;
                stream.write_all(format!(":{} 004 {} {} {} {} {}\r\n", nickname, nickname, servername, version, user_modes, channel_modes).as_bytes())?;
                stream.write_all(format!(":{} 005 {} :Try server {}, port {}\r\n", nickname, nickname, server_name, port_number).as_bytes())?;
                stream.write_all(format!(":{} 375 {} :- {} Message of the day -\r\n", nickname, nickname, servername).as_bytes())?;
                stream.write_all(format!(":{} 372 {} :- {}\r\n", nickname, nickname, motd_text).as_bytes())?;
                stream.write_all(format!(":{} 376 {} :End of /MOTD command.\r\n", nickname, nickname).as_bytes())?;
            }
            Command::User(username) => {
                stream.write_all(format!("USER {}\r\n", username).as_bytes())?;
            }
            Command::Quit => {
                stream.write_all("Goodbye!\r\n".as_bytes())?;
                return Ok(());
            }
            Command::Join(channel) => {
                let mut channels = channels.lock().unwrap();
                channels.entry(channel.clone()).or_insert_with(Vec::new).push(nickname.clone());
                stream.write_all(format!(":{} JOIN {}\r\n", nickname, channel).as_bytes())?;
                if let Some(members) = channels.get(&channel) {
                    for member in members {
                        stream.write_all(format!(":{} 353 {} = {} :{}\r\n", nickname, nickname, channel, member).as_bytes())?;
                    }
                    stream.write_all(format!(":{} 366 {} {} :End of /NAMES list.\r\n", nickname, nickname, channel).as_bytes())?;
                }
                // Make the user an operator when they join a channel
                let mut ops = ops.lock().unwrap();
                ops.insert(channel.clone(), nickname.clone());
                stream.write_all(format!(":{} MODE {} +o {}\r\n", nickname, channel, nickname).as_bytes())?;
            }
            Command::Part(channel) => {
                let mut channels = channels.lock().unwrap();
                if let Some(members) = channels.get_mut(&channel) {
                    members.retain(|nick| nick != &nickname);
                    stream.write_all(format!(":{} PART {}\r\n", nickname, channel).as_bytes())?;
                } else {
                    stream.write_all(format!("You're not in channel {}\r\n", channel).as_bytes())?;
                }
            }
            Command::PrivMsg(target, message) => {
                stream.write_all(format!("PRIVMSG {} :{}\r\n", target, message).as_bytes())?;
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
            Command::Names(channel) => {
                let channels = channels.lock().unwrap();
                let ops = ops.lock().unwrap();
                if let Some(members) = channels.get(&channel) {
                    let names: Vec<String> = members.iter().map(|member| {
                        if ops.get(&channel) == Some(member) {
                            format!("@{}", member)
                        } else {
                            member.clone()
                        }
                    }).collect();
                    let names = names.join(" ");
                    stream.write_all(format!(":{} 353 {} = {} :{}\r\n", nickname, nickname, channel, names).as_bytes())?;
                    stream.write_all(format!(":{} 366 {} {} :End of /NAMES list.\r\n", nickname, nickname, channel).as_bytes())?;
                } else {
                    stream.write_all(format!("No such channel: {}\r\n", channel).as_bytes())?;
                }
            }
            Command::Mode(channel, mode, user) => {
                println!("Handling MODE command"); // Debug print
                let mut ops = ops.lock().unwrap();
                if let Some(op) = ops.get(&channel) {
                    if *op == nickname {
                        if mode == "+o" {
                            ops.insert(channel.clone(), user.clone());
                            stream.write_all(format!(":{} MODE {} +o {}\r\n", nickname, channel, user).as_bytes())?;
                        } else if mode == "-o" {
                            if user == *op {
                                ops.remove(&channel);
                                stream.write_all(format!(":{} MODE {} -o {}\r\n", nickname, channel, user).as_bytes())?;
                            } else {
                                stream.write_all(format!("You're not an operator of channel {}\r\n", channel).as_bytes())?;
                            }
                        }
                    } else {
                        stream.write_all(format!("You're not an operator of channel {}\r\n", channel).as_bytes())?;
                    }
                } else {
                    stream.write_all(format!("No operator for channel {}\r\n", channel).as_bytes())?;
                }
            }
            Command::Who(channel) => {
                // Handle WHO command
                let channels = channels.lock().unwrap();
                if let Some(members) = channels.get(&channel) {
                    let names: Vec<String> = members.iter().map(|member| member.clone()).collect();
                    let names = names.join(" ");
                    stream.write_all(format!(":{} 352 {} {} {} {} {} H :0 {}\r\n", nickname, nickname, channel, nickname, "localhost", "localhost", names).as_bytes())?;
                    stream.write_all(format!(":{} 315 {} {} :End of /WHO list.\r\n", nickname, nickname, channel).as_bytes())?;
                } else {
                    stream.write_all(format!("No such channel: {}\r\n", channel).as_bytes())?;
                }
            }
            Command::Unknown(cmd) => {
                println!("Unknown command: {}", cmd); // Debug print
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
    for stream in listener.incoming() {
        match stream {
            Ok(stream) => {
                let channels = Arc::clone(&channels);
                let ops = Arc::clone(&ops);
                thread::spawn(move || {
                    if let Err(e) = handle_client(stream, channels, ops) {
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
