import { expect, test } from "vitest"
import {
    fs,
} from "./challenge.js"

test("test function-state", () => {
    expect(fs()).toEqual(42);
    expect(fs()).toEqual(420);
    expect(fs()).toEqual(69);
    expect(fs()).toEqual(69);
    expect(fs()).toEqual(420);
    expect(fs()).toEqual(42);
    expect(fs() && fs.foo).toEqual(7);
    expect(fs() && fs.foo).toEqual(8);
    expect(fs() && fs.foo).toEqual(9);
    expect(fs()).toEqual(null);
    expect(fs()).toEqual(0xDEADBEA7)
    expect((fs.foo = fs()).foo).toEqual([69, "nice"])
});



