# Text Completion/Editing/Auto-correction Tool

This is a Go program designed as part of the curriculum of [01Talent/Zone01 Oujda](https://github.com/01-edu/public/tree/master/subjects/go-reloaded). It is a text completion, editing, and auto-correction tool that performs various modifications on input text according to specific instructions. The tool implements transformations such as converting text to uppercase or lowercase, capitalizing words, converting hexadecimal or binary numbers to decimal, formatting punctuation marks, and handling articles before vowels.

## Features

- **Transformation Actions**: The program recognizes transformation actions such as `(up)`, `(low)`, `(cap)`, `(hex)`, and `(bin)`, and performs the corresponding modifications to the text.
- **Numbered Transformations**: It supports transformations with a specified number of words, like `(up, 2)` to convert the two previous words to uppercase.
- **Punctuation Formatting**: Punctuation marks like `.`, `,`, `!`, `?`, `:`, and `;` are formatted according to the specified rules. Additionally, special cases like `...` and `!?` are handled appropriately.
- **Quotation Handling**: The program correctly handles quotation marks, ensuring they are placed correctly around words or phrases.
- **Article Handling**: It automatically changes the article "a" to "an" if followed by a word starting with a vowel or "h".

## Usage

To use the program, follow these steps:

1. Ensure you have Go installed on your system.
2. Clone this repository to your local machine.
3. Navigate to the directory containing the code.
4. Prepare your input text file containing the text to be modified.
5. Run the program using the command-line interface, providing the input and output file paths as arguments.

Example usage:

```bash
go run . input.txt output.txt
```
Replace `input.txt` with the path to your input file and `output.txt` with the desired output file path.

### Example

Suppose you have an input file `input.txt` with the following content:
```plaintext
Simply add 42 `(hex)` and 10 `(bin)` and you will see the result is 68.
```
```bash
./go-reloaded input.txt output.txt
```

```plaintext
Simply add `66` and <span style= "color : blue;">2</span>` and you will see the result is 68.
```



