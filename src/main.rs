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
    PrivMsg(String, String),
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
        Some(&"PRIVMSG") if parts.len() >= 3 => {
            let target = parts[1].to_string();
            let message = parts[2..].join(" ");
            Command::PrivMsg(target, message)
        }
        _ => Command::Unknown(line.to_string()),
    }
}

fn handle_client(mut stream: TcpStream, channels: Arc<Mutex<HashMap<String, Vec<String>>>>) -> std::io::Result<()> {
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
                stream.write_all(format!("NICK {}\r\n", nick).as_bytes())?;
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
            }
            Command::PrivMsg(target, message) => {
                stream.write_all(format!("PRIVMSG {} :{}\r\n", target, message).as_bytes())?;
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
    for stream in listener.incoming() {
        match stream {
            Ok(stream) => {
                let channels = Arc::clone(&channels);
                thread::spawn(move || {
                    if let Err(e) = handle_client(stream, channels) {
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