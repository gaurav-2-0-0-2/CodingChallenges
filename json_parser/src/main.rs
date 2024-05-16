// JSON --> Javascript Object Notation
/*  
 
    {
      key1: value1,
      key2: value2,
    }
 
 */

use std::fs;
//use std::io::prelude::*;
//use std::path::Path;
use std::env;


fn lex(string: &String){
   let mut vec = Vec::new();
   for i in string.chars(){
       vec.push(i);
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
