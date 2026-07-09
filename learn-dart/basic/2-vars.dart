void main() {
  // 1. Declare a variable

  // String
  var name = 'John Doe';
  String name2 = 'Jane Doe';
  print('Name: $name');
  print('Name2: $name2');

  // int
  var age = 24;
  int age2 = 25;
  print('Age: $age');
  print('Age2: $age2');

  // dynamic
  // is a special type in Dart that turns off compile-time type checking for a variable — it lets
  // the variable hold a value of any type, and lets you call any method or property on it, with those
  // checks deferred to runtime instead of caught by the analyzer
  dynamic dynamicValue = 42;
  print('Dynamic Value: $dynamicValue');
  // value = "hello"; // fine, no error
  // value = [1, 2, 3]; // also fine
  // value
  //     .someRandomMethod(); // compiles fine — but crashes at runtime if it doesn't exist

  // 2. Declare final and constant

  // final
  // means "this variable can only be set once" — but
  // the value can be determined at runtime.
  final name3 = "Alice";
  final now = DateTime.now(); // fine — computed at runtime
  final list = [1, 2, 3];
  list.add(4); // fine! the list itself is mutable, just can't reassign `list`
  // list = [5, 6, 7];             // ERROR — can't reassign

  // const
  // means "this variable is a compile-time constant" — the value
  // must be known at compile time, and the value is deeply immutable.
  const name4 = "Alice";
  const pi = 3.14159;
  // const now = DateTime.now();   // ERROR — not known at compile time
  const list2 = [1, 2, 3];
  // list2.add(4);                  // ERROR — const collections are frozen, can't mutate contents either
}
