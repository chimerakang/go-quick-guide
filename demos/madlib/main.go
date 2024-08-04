/*
*
Madlibs program is a software that takes in a users inputs and places
these inputs into a String, then prints it out

@author: Uchenna Nwoke
@version: 1.0
*/
package main

import (
	"bufio"
	"fmt"
	"os"
)

/*
*
This is the main method, which creates a scanner object to accept input from
standard in. These inputs are stored in variables and used in a setence.
*/
func main() {
	//  scanner object to read from Standard input
	scanner := bufio.NewScanner(os.Stdin)

	/*
		Prompts a question to standard out
		Scans for any entries made in the standard in (command line/terminal)
		Passes those values as string and stores in the variable
	*/
	fmt.Print("Enter an adjective: ")
	scanner.Scan()
	adj := scanner.Text()

	fmt.Print("Enter a verb: ")
	scanner.Scan()
	verb1 := scanner.Text()

	fmt.Print("Enter another verb: ")
	scanner.Scan()
	verb2 := scanner.Text()

	fmt.Print("Enter a famous person: ")
	scanner.Scan()
	famousPerson := scanner.Text()

	// Using the fmt.Sprintf to format the string
	madlibs := fmt.Sprintf("Computer pramming is so %s! It makes me so excited all the time because \n"+
		" I love to %s. Stay hydrated and %s like you are %s", adj, verb1, verb2, famousPerson)

	// Prints madlibs to standard out
	fmt.Printf(madlibs)
}
