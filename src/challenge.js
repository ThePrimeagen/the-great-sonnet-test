class Vector {
    constructor(x, y) {
        this.x = x;
        this.y = y;
    }

    mul(vector) {
        return new Vector(this.x * vector.x, this.y * vector.y);
    }

    equals(vector) {
        return this.x === vector.x && this.y === vector.y;
    }
}

function run(i) {
    if (i === 13) {
        return new Vector(3, 4);
    } else if (i === 42) {
        return () => 1337;
    } else if (i % 10 === 0) {
        return 69;
    } else {
        return 42;
    }
}

export { Vector, run };
