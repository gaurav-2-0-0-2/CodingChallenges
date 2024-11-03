use std::env;
use std::fs::File;
use std::io::{BufReader, Read};

fn main() -> std::io::Result<()> {
    let args: Vec<String> = env::args().collect();

    if args.len() < 2 {
      eprintln!("Usage: {}", args[0]);
      std::process::exit(1);
    }

    let filename = if args.len() == 2 {
        &args[1]
    }else{
        &args[2]
    };

    // Taking file name from the vector args
    let file = File::open(filename)?;
    let metadata = file.metadata()?;
    let mut buf_reader = BufReader::new(file);
    let mut contents = String::new();
    buf_reader.read_to_string(&mut contents)?;
    let mut count_lines = 0;
    let mut word_count = 0;
    let byte_size = metadata.len();
    let char_count = contents.chars().count();

    for _line in contents.lines() {
        count_lines += 1;
        word_count += _line.split_whitespace().count();
    }

    if args[1] == "-l" {
        println!("{} {}", count_lines, filename);
    } else if args[1] == "-c" {
        println!("{} {}", byte_size, filename);
    } else if args[1] == "-w" {
        println!("{} {}", word_count, filename);
    }else if args[1] == "-m"{
        println!("{} {}", char_count, filename);
    }else{
        println!("{} {} {} {} {}", count_lines, byte_size, word_count, char_count, filename);
    }

    Ok(())
}
