let counter = 0;
let fooCounter = 7;
const sequence = [42, 420, 69, 69, 420, 42];

function fs(input) {
    if (input !== undefined) {
        return input === 6 ? 42 : input * 2;
    }

    if (counter < sequence.length) {
        return sequence[counter++];
    }

    if (counter === sequence.length) {
        counter++;
        return true;
    }

    if (counter === sequence.length + 1) {
        counter++;
        return true;
    }

    if (counter === sequence.length + 2) {
        counter++;
        return true;
    }

    if (counter === sequence.length + 3) {
        counter++;
        return null;
    }

    if (counter === sequence.length + 4) {
        counter++;
        return 0xDEADBEA7;
    }

    if (counter === sequence.length + 5) {
        counter++;
        return { foo: [69, "nice"] };
    }

    return 0;
}

Object.defineProperty(fs, 'foo', {
    get: function() {
        return fooCounter++;
    },
    set: function(value) {
        if (typeof value === 'object' && value.foo) {
            this.foo = value.foo;
        }
    }
});

export { fs };
