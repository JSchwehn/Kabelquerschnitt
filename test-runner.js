// Test runner for Node.js to validate test logic
const COPPER_RESISTIVITY_20C = 0.0175;
const ALUMINUM_RESISTIVITY_20C = 0.0283;
const COPPER_TEMP_COEFFICIENT = 0.00393;
const ALUMINUM_TEMP_COEFFICIENT = 0.00403;
const REFERENCE_TEMP = 25.0;

const materials = {
    copper: {
        nameKey: "materialCopper",
        resistivity20C: COPPER_RESISTIVITY_20C,
        tempCoefficient: COPPER_TEMP_COEFFICIENT
    },
    aluminum: {
        nameKey: "materialAluminum",
        resistivity20C: ALUMINUM_RESISTIVITY_20C,
        tempCoefficient: ALUMINUM_TEMP_COEFFICIENT
    }
};

const installationAdjustments = {
    air: 0.0,
    conduit: 10.0,
    isolated: 20.0
};

const standardMetricSizes = [0.5, 0.75, 1.0, 1.5, 2.5, 4.0, 6.0, 10.0, 16.0, 25.0, 35.0, 50.0, 70.0, 95.0, 120.0, 150.0, 185.0, 240.0];

const awgSizes = [
    { label: "18", area: 0.823 },
    { label: "16", area: 1.309 },
    { label: "14", area: 2.081 },
    { label: "12", area: 3.309 },
    { label: "10", area: 5.261 },
    { label: "8", area: 8.367 },
    { label: "6", area: 13.30 },
    { label: "4", area: 21.15 },
    { label: "2", area: 33.62 },
    { label: "1", area: 42.41 },
    { label: "1/0", area: 53.49 },
    { label: "2/0", area: 67.43 },
    { label: "3/0", area: 85.01 },
    { label: "4/0", area: 107.2 }
];

const wireTypes = {
    generic: { name: "Generic", maxTemp: 90.0 },
    flry: { name: "FLRY", maxTemp: 105.0 },
    pvc: { name: "PVC", maxTemp: 70.0 }
};

function fahrenheitToCelsius(f) {
    return (f - 32) * 5 / 9;
}

function celsiusToFahrenheit(c) {
    return c * 9 / 5 + 32;
}

function calculateResistivityAtTemp(material, tempCelsius) {
    return material.resistivity20C * (1 + material.tempCoefficient * (tempCelsius - REFERENCE_TEMP));
}

function calculateEffectiveTemp(ambientTempCelsius, installation) {
    return ambientTempCelsius + installationAdjustments[installation];
}

function validateWireTemperature(effectiveTempCelsius, wireType) {
    if (effectiveTempCelsius > wireType.maxTemp) {
        return { isValid: false, message: `Temperature exceeds max` };
    }
    if (effectiveTempCelsius > wireType.maxTemp * 0.9) {
        return { isValid: true, message: `Temperature close to max` };
    }
    return { isValid: true, message: "" };
}

function calculateCableArea(voltage, current, length, maxVoltageDropPercent, material, roundTrip, ambientTempCelsius, installation) {
    const maxVoltageDrop = voltage * (maxVoltageDropPercent / 100.0);
    const distanceFactor = roundTrip ? 2.0 : 1.0;
    const effectiveTemp = calculateEffectiveTemp(ambientTempCelsius, installation);
    const resistivity = calculateResistivityAtTemp(material, effectiveTemp);
    return (current * resistivity * length * distanceFactor) / maxVoltageDrop;
}

function areaToDiameter(area) {
    return 2 * Math.sqrt(area / Math.PI);
}

function findClosestMetricSize(requiredArea) {
    for (const size of standardMetricSizes) {
        if (size >= requiredArea) {
            return { size: size, diff: size - requiredArea };
        }
    }
    const largestSize = standardMetricSizes[standardMetricSizes.length - 1];
    return { size: largestSize, diff: requiredArea - largestSize };
}

function findClosestAWG(requiredArea) {
    for (const awg of awgSizes) {
        if (awg.area >= requiredArea) {
            return { label: awg.label, area: awg.area, diff: awg.area - requiredArea };
        }
    }
    const largestAWG = awgSizes[awgSizes.length - 1];
    return { label: largestAWG.label, area: largestAWG.area, diff: requiredArea - largestAWG.area };
}

// Test runner
let passed = 0;
let failed = 0;

function assert(condition, message) {
    if (!condition) {
        throw new Error(message || 'Assertion failed');
    }
}

function assertEqual(actual, expected, message) {
    const tolerance = 0.0001;
    const diff = Math.abs(actual - expected);
    if (diff > tolerance) {
        throw new Error(message || `Expected ${expected}, got ${actual} (diff: ${diff})`);
    }
}

function assertClose(actual, expected, tolerance, message) {
    const diff = Math.abs(actual - expected);
    if (diff > tolerance) {
        throw new Error(message || `Expected ${expected} ± ${tolerance}, got ${actual} (diff: ${diff})`);
    }
}

function test(name, fn) {
    try {
        fn();
        console.log(`✅ ${name}`);
        passed++;
    } catch (error) {
        console.log(`❌ ${name}`);
        console.log(`   Error: ${error.message}`);
        failed++;
    }
}

// Run tests
console.log('Running tests...\n');

test('fahrenheitToCelsius - converts 32°F to 0°C', () => {
    assertEqual(fahrenheitToCelsius(32), 0);
});

test('fahrenheitToCelsius - converts 212°F to 100°C', () => {
    assertEqual(fahrenheitToCelsius(212), 100);
});

test('fahrenheitToCelsius - converts 68°F to 20°C', () => {
    assertEqual(fahrenheitToCelsius(68), 20);
});

test('celsiusToFahrenheit - converts 0°C to 32°F', () => {
    assertEqual(celsiusToFahrenheit(0), 32);
});

test('celsiusToFahrenheit - converts 100°C to 212°F', () => {
    assertEqual(celsiusToFahrenheit(100), 212);
});

test('celsiusToFahrenheit - converts 20°C to 68°F', () => {
    assertEqual(celsiusToFahrenheit(20), 68);
});

test('calculateResistivityAtTemp - copper at 25°C (reference temp)', () => {
    const resistivity = calculateResistivityAtTemp(materials.copper, 25.0);
    assertEqual(resistivity, COPPER_RESISTIVITY_20C);
});

test('calculateResistivityAtTemp - copper at 40°C increases', () => {
    const resistivity25 = calculateResistivityAtTemp(materials.copper, 25.0);
    const resistivity40 = calculateResistivityAtTemp(materials.copper, 40.0);
    assert(resistivity40 > resistivity25, 'Resistivity should increase with temperature');
});

test('calculateResistivityAtTemp - aluminum at 25°C (reference temp)', () => {
    const resistivity = calculateResistivityAtTemp(materials.aluminum, 25.0);
    assertEqual(resistivity, ALUMINUM_RESISTIVITY_20C);
});

test('calculateEffectiveTemp - air installation (no adjustment)', () => {
    const temp = calculateEffectiveTemp(25.0, 'air');
    assertEqual(temp, 25.0);
});

test('calculateEffectiveTemp - conduit installation (+10°C)', () => {
    const temp = calculateEffectiveTemp(25.0, 'conduit');
    assertEqual(temp, 35.0);
});

test('calculateEffectiveTemp - isolated installation (+20°C)', () => {
    const temp = calculateEffectiveTemp(25.0, 'isolated');
    assertEqual(temp, 45.0);
});

test('validateWireTemperature - valid temperature', () => {
    const result = validateWireTemperature(50.0, wireTypes.generic);
    assert(result.isValid === true, 'Should be valid');
    assert(result.message === '', 'Should have no warning');
});

test('validateWireTemperature - temperature exceeds max', () => {
    const result = validateWireTemperature(100.0, wireTypes.generic);
    assert(result.isValid === false, 'Should be invalid');
    assert(result.message !== '', 'Should have warning message');
});

test('validateWireTemperature - temperature close to max (caution)', () => {
    const result = validateWireTemperature(85.0, wireTypes.generic);
    assert(result.isValid === true, 'Should be valid but with caution');
    assert(result.message !== '', 'Should have caution message');
});

test('areaToDiameter - calculates diameter correctly', () => {
    const area = Math.PI;
    const diameter = areaToDiameter(area);
    assertClose(diameter, 2.0, 0.0001);
});

test('areaToDiameter - round trip verification', () => {
    const testArea = 10.0;
    const diameter = areaToDiameter(testArea);
    const radius = diameter / 2;
    const calculatedArea = Math.PI * radius * radius;
    assertClose(calculatedArea, testArea, 0.0001);
});

test('findClosestMetricSize - rounds up correctly', () => {
    const result = findClosestMetricSize(3.5);
    assert(result.size >= 3.5, 'Should round up');
    assertEqual(result.size, 4.0);
    assert(result.diff >= 0, 'Diff should be positive when rounding up');
});

test('findClosestMetricSize - exact match', () => {
    const result = findClosestMetricSize(4.0);
    assertEqual(result.size, 4.0);
    assertEqual(result.diff, 0.0);
});

test('findClosestMetricSize - very small area', () => {
    const result = findClosestMetricSize(0.3);
    assertEqual(result.size, 0.5);
});

test('findClosestMetricSize - exceeds all sizes', () => {
    const result = findClosestMetricSize(500.0);
    assertEqual(result.size, 240.0);
    assert(result.diff > 0, 'Diff should be positive (requiredArea - largestSize)');
});

test('findClosestAWG - rounds up correctly', () => {
    const result = findClosestAWG(3.5);
    assert(result.area >= 3.5, 'Should round up');
    assertEqual(result.label, '10');
    assertEqual(result.area, 5.261);
});

test('findClosestAWG - exact match', () => {
    const result = findClosestAWG(5.261);
    assertEqual(result.label, '10');
    assertEqual(result.area, 5.261);
});

test('findClosestAWG - very small area', () => {
    const result = findClosestAWG(0.5);
    assertEqual(result.label, '18');
});

test('findClosestAWG - exceeds all sizes', () => {
    const result = findClosestAWG(200.0);
    assertEqual(result.label, '4/0');
    assert(result.diff > 0, 'Diff should be positive (requiredArea - largestAWG.area)');
});

test('calculateCableArea - 12V, 10A, 5m, 3%, copper, one-way, 20°C, air', () => {
    const area = calculateCableArea(12, 10, 5, 3, materials.copper, false, 20, 'air');
    assertClose(area, 2.38, 0.1);
});

test('calculateCableArea - round trip doubles area requirement', () => {
    const areaOneWay = calculateCableArea(12, 10, 5, 3, materials.copper, false, 20, 'air');
    const areaRoundTrip = calculateCableArea(12, 10, 5, 3, materials.copper, true, 20, 'air');
    assertClose(areaRoundTrip, areaOneWay * 2, 0.01);
});

test('calculateCableArea - higher temperature requires larger area', () => {
    const area20 = calculateCableArea(12, 10, 5, 3, materials.copper, false, 20, 'air');
    const area40 = calculateCableArea(12, 10, 5, 3, materials.copper, false, 40, 'air');
    assert(area40 > area20, 'Higher temperature should require larger area');
});

test('calculateCableArea - aluminum requires larger area than copper', () => {
    const areaCopper = calculateCableArea(12, 10, 5, 3, materials.copper, false, 20, 'air');
    const areaAluminum = calculateCableArea(12, 10, 5, 3, materials.aluminum, false, 20, 'air');
    assert(areaAluminum > areaCopper, 'Aluminum should require larger area');
});

test('calculateCableArea - conduit installation increases area requirement', () => {
    const areaAir = calculateCableArea(12, 10, 5, 3, materials.copper, false, 20, 'air');
    const areaConduit = calculateCableArea(12, 10, 5, 3, materials.copper, false, 20, 'conduit');
    assert(areaConduit > areaAir, 'Conduit installation should require larger area');
});

test('calculateCableArea - higher voltage drop allows smaller area', () => {
    const area3 = calculateCableArea(12, 10, 5, 3, materials.copper, false, 20, 'air');
    const area5 = calculateCableArea(12, 10, 5, 5, materials.copper, false, 20, 'air');
    assert(area5 < area3, 'Higher voltage drop should allow smaller area');
});

// Summary
console.log(`\n=== Test Summary ===`);
console.log(`Total: ${passed + failed} tests`);
console.log(`Passed: ${passed}`);
console.log(`Failed: ${failed}`);
console.log(`Pass rate: ${((passed / (passed + failed)) * 100).toFixed(1)}%`);

process.exit(failed > 0 ? 1 : 0);

