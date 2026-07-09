void main() {
  // Functions

  // Named optional parameters use {}
  // Called like: myFunc("John", name2: "Bill")
  String myFunc(String name1, {String name2 = "friends"}) {
    return "Hello $name1 and $name2";
  }

  var thing = myFunc("John", name2: "Bill");
  print(thing);

  // Positional optional parameters use []
  // No name needed when calling — just pass it in order
  // Called like: myFunc2("John") or myFunc2("John", "Bill")
  String myFunc2(String name1, [String name2 = "friends"]) {
    return "Hello $name1 and $name2";
  }

  var thing2 = myFunc2("John", "Bill");
  print(thing2);

  var thing3 = myFunc2("John"); // uses default "friends"
  print(thing3);
}
