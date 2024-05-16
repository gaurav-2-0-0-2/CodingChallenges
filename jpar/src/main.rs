use std::fs;
use std::env;

fn lex(input_string: &String){
    let mut vec = Vec::new();
    for c in input_string.chars(){
        vec.push(c);
    }
    println!("{:?}", vec);
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let file_path = &args[1];
    let contents = fs::read_to_string(file_path)
        .expect("Should have been able to read the file");
    lex(&contents);
}
