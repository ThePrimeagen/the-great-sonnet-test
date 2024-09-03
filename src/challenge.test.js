import { expect, test } from "vitest"
import * as challenge from "./challenge.js"

test("The Great Sonnet Challenge", () => {
    for (let i = 0; i < 100; ++i) {
        const out = challenge.run(i);
        if (i == 13) {
            expect(out.mul(new challenge.Vector(5, 2))).toEqual(new challenge.Vector(15, 8));
        } else if (i == 42) {
            expect(out()).toEqual(1337);
        } else if (i % 10 == 0) {
            expect(out).toEqual(69);
        } else {
            expect(out).toEqual(42);
        }
    }
})

test("Another Great Test", () => {
    const vec = new challenge.Vector(2, 4);

    expect(vec.x).toEqual(2);
    expect(vec.y).toEqual(4);
})
