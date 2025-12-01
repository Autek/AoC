use std::fs;

struct Left {
    val: i32,
}

struct Right {
    val: i32,
}

enum Code {
    Left(Left),
    Right(Right),
}

fn print_type_of<T>(_: &T) {
    println!("\ntype : {}\n", std::any::type_name::<T>());
}

fn main() {
    let file_path = "data.txt";
    let contents_string = fs::read_to_string(file_path).unwrap();

    let (final_acc, zero_hits) = contents_string
        .lines()
        .filter(|line| !line.trim().is_empty())
        .map(|line| {
            let (kind, rest) = line.split_at(1);
            let val: i32 = rest.parse().unwrap();

            match kind {
                "R" => Code::Right(Right { val }),
                "L" => Code::Left(Left { val }),
                _ => panic!("Unknown kind: {kind}"),
            }
        })
        .fold(
            (50, 0),
            |(acc, count), code| {
                let new_acc = match code {
                    Code::Right(Right { val }) => (acc + val) % 100,
                    Code::Left(Left { val }) => (acc - val + 100) % 100,
                };

                let new_count = if new_acc == 0 { count + 1 } else { count };

                (new_acc, new_count)
            },
        );

    println!("final_acc = {final_acc}, zero_hits = {zero_hits}");
}

