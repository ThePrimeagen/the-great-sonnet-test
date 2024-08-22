import { expect, test } from "vitest"
import {
    challenge,
    Vector,
} from "./challenge.js"

test("runs the challenge function 100 times", () => {
    for (let i = 0; i < 100; ++i) {
        const out = challenge();

        if (i % 10 == 0) {
            expect(out).toEqual(69);
        } else if (i == 69) {
            expect(out()).toEqual("nice");
        } else if (i == 42) {
            expect(out.mul(new Vector([2, 3, 4]))).toEqual(new Vector([
                42,
                69,
                420
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

