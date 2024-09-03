import { expect, test } from "vitest"
import challenge from "./challenge.js"

test("runs the challenge run function 100 times", () => {
    for (let i = 0; i < 100; ++i) {
        const out = challenge.run();
        expect(out).toEqual(42);
    }
})

