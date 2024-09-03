export function run() {
    return 42;
}

function calculateTotal(items) {
    return items.reduce((total, item) => total + item.price * item.quantity, 0);
}

function applyDiscount(total, discountPercent) {
    return total * (1 - discountPercent / 100);
}

function calculateTax(total, taxRate) {
    return total * (taxRate / 100);
}

function generateInvoice(items, discountPercent, taxRate) {
    const subtotal = calculateTotal(items);
    const discountedTotal = applyDiscount(subtotal, discountPercent);
    const tax = calculateTax(discountedTotal, taxRate);
    const total = discountedTotal + tax;

    return {
        subtotal: Number(subtotal.toFixed(2)),
        discount: Number((subtotal - discountedTotal).toFixed(2)),
        tax: Number(tax.toFixed(2)),
        total: Number(total.toFixed(2))
    };
}

