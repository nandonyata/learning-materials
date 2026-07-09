void main() {
  // Declare a variable

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
}
