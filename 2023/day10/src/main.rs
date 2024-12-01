use std::fs::File;
use std::io::BufRead;

fn main() -> Result<(), std::io::Error> {
    let fname = std::env::args().nth(1).unwrap();
    let file = File::open(fname)?;
    let reader = std::io::BufReader::new(file);

    // TODO: Solve a problem!
    Ok(())
}

mod test {
    use super::*;

    #[test]
    fn test_something() {
    }
}
