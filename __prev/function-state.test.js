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
    expect(fs()).toEqual(0xDEADBEA7);
    expect((fs.foo = fs()).foo).toEqual([69, "nice"]);
    expect(fs(0)).toEqual(0);
    expect(fs(1)).toEqual(2);
    expect(fs(2)).toEqual(4);
    expect(fs(3)).toEqual(6);
    expect(fs(4)).toEqual(8);
    expect(fs(5)).toEqual(10);
    expect(fs(6)).toEqual(42);
});



