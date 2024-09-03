import { expect, test } from "vitest"
import {
    challenge,
    number2,
    Vector,
} from "./challenge.js"

function isPrime(num) {
    if (num == 1 || num == 0) {
        return false
    }

    if (num < 4) {
        return true;
    }

    if (num & 0x1 == 0) {
        return false;
    }

    for (let i = 5; i < Math.sqrt(num); i += 2) {
        if (num % i == 0) {
            return false;
        }
    }
    return true;
}

test("runs the challenge function 100 times", () => {
    let _42Vec = null
    for (let i = 0; i < 100; ++i) {
        const out = challenge();

        if (isPrime(i)) {
            expect(out).toEqual("prime");
        } else if (i % 10 == 0) {
            expect(out).toEqual(69);
        } else if (i == 69) {
            expect(out()).toEqual("nice");
        } else if (i == 42) {
            expect(out.mul(new Vector([2, 3, 4]))).toEqual(new Vector([
                42,
                69,
                420
            ]));
        } else if (i == 43) {
            expect(out.mul(new Vector([0, 0, 0]))).toEqual(new Vector([
                0,
                0,
                0
            ]));
        } else {
            expect(out).toEqual(42);
        }
    }
})

test("perform vector operations", () => {
    const vec = new Vector([3, 5, 7]);
    const vec2 = new Vector([66, 64, 62]);

    expect(vec.add(vec2)).toEqual(new Vector([69, 69, 69]));
})

