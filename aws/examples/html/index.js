function helloWorld () {
  console.debug('~~~ Hello, world!');

  let num = window.prompt('Please enter a number:', '42');
  num = parseInt(num, 10);
  if (isNaN(num) || num <= 0) {
    num = 42;
  }

  console.log(`~~~ num: ${num}`);
}
