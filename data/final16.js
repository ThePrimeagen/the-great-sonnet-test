export function challenge() {
    let count = 0;
    return function() {
        if (count % 10 === 0) {
            count++;
            return 69;
        } else if (count === 69) {
            count++;
            return () => "nice";
        } else {
            count++;
            return 42;
        }
    };
}

export class Vector {
    constructor(components) {
        this.components = components;
    }

    add(other) {
        if (this.components.length !== other.components.length) {
            throw new Error("Vectors must have the same dimension");
        }
        const result = this.components.map((value, index) => value + other.components[index]);
        return new Vector(result);
    }
}