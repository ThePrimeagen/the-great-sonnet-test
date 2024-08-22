import { expect, test } from "vitest"
import {
    challenge
} from "./challenge.js"

test("runs the challenge function 100 times", () => {
    for (let i = 0; i < 100; ++i) {
        const out = challenge(i);
        const expected = i % 10 == 0 ? 69 : 42;
        expect(out).toEqual(expected);
    }
})

