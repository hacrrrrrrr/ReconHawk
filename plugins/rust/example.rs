use std::env;

fn main() {
    let target = env::args().nth(1).unwrap_or_default();
    println!(r#"{{"plugin":"rust-example","target":{:?},"type":"enrichment"}}"#, target);
}
