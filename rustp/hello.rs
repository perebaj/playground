fn main() {
    println!("Hello, world!");
    // replace using {}
    println!("{} days", 31);
    // Positional arguments
    println!("{0}, this is {1}. {1}, this is {0}", "Alice", "Bob");
    // Named arguments
    println!("{subject} {verb} {object}",
            object="the lazy dog",
            subject="the quick brown fox",
            verb="jumps over");
    
    // Special formatting
    println!("{} of {:b} people know binary, the other half doesn't", 1, 2);
    // Base 2 = binary {:b}
    // Base 8 = octal {:o}
    // Base 16 = hex {:x}
    // Base 16 = hex {:X}
    println!("Base 16 {:x} {:X}", 255, 255);

    // Error when the positional argument is not found
    // println!("My name is {0}, {1} {0}", "Bond");

    let number: f64 = 1.0;
    let width: usize = 2;

    println!("{number:>0width$}", number=number, width=width);

    // struct UnPrintable(i32);
    #[derive(Debug)]
    struct DebugPrintable(i32);

    // println!("This struct `{:?}` won't print...", UnPrintable(3));
    println!("This struct `{:?}` will print...", DebugPrintable(3));

    //Pretty print
    println!("Now {:#?} will print!", DebugPrintable(3));
}


use std::fmt;

struct Structure(i32);

impl fmt::Display for Structure {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        // Write strictly the first element into the supplied output
        // stream: f
        write!(f, "{}", self.0)
    }
}
