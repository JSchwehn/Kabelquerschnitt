// Language and translations
// Translations are loaded from translations.js
let currentLanguage = localStorage.getItem('language') || 'de';
// Translation function
function t(key, params = {}) {
    let text = translations[currentLanguage][key] || translations['de'][key] || key;
    Object.keys(params).forEach(param => {
        text = text.replace(`{${param}}`, params[param]);
    });
    return text;
}

function changeLanguage(lang) {
    currentLanguage = lang;
    localStorage.setItem('language', lang);
    document.documentElement.lang = lang;
    updateTranslations();
}

function updateTranslations() {
    // Update all elements with data-i18n attribute
    document.querySelectorAll('[data-i18n]').forEach(el => {
        const key = el.getAttribute('data-i18n');
        if (translations[currentLanguage][key]) {
            el.textContent = t(key);
        }
    });

    // Update select options
    document.getElementById('roundTrip').innerHTML = `
        <option value="false">${t('oneWay')}</option>
        <option value="true">${t('roundTrip')}</option>
    `;

    document.getElementById('material').innerHTML = `
        <option value="copper">${t('materialCopper')}</option>
        <option value="aluminum">${t('materialAluminum')}</option>
    `;

    document.getElementById('installation').innerHTML = `
        <option value="air">${t('installationAir')}</option>
        <option value="conduit">${t('installationConduit')} (+10°C)</option>
        <option value="isolated">${t('installationIsolated')} (+20°C)</option>
    `;

    // Update button
    const btn = document.getElementById('calculateBtn');
    if (btn) btn.textContent = t('calculate');
}

// Constants
const MAX_VOLTAGE = 60.0; // Maximum system voltage (V)
const COPPER_RESISTIVITY_20C = 0.0175; // Ω·mm²/m
const ALUMINUM_RESISTIVITY_20C = 0.0283; // Ω·mm²/m
const COPPER_TEMP_COEFFICIENT = 0.00393; // per °C
const ALUMINUM_TEMP_COEFFICIENT = 0.00403; // per °C
const REFERENCE_TEMP = 20.0; // °C

// Initialize language on page load
document.addEventListener('DOMContentLoaded', function () {
    // Set language selector
    document.getElementById('language').value = currentLanguage;
    document.documentElement.lang = currentLanguage;

    // Set max voltage attribute
    document.getElementById('voltage').setAttribute('max', MAX_VOLTAGE);

    // Update translations
    updateTranslations();

    // Allow Enter key to trigger calculation
    const inputs = document.querySelectorAll('input, select');
    inputs.forEach(input => {
        input.addEventListener('keypress', function (e) {
            if (e.key === 'Enter') {
                calculate();
            }
        });
    });
});

// Material densities (kg/m³)
const COPPER_DENSITY = 8960; // kg/m³
const ALUMINUM_DENSITY = 2700; // kg/m³

// Weight per meter per mm² (g/m per mm²)
const COPPER_WEIGHT_PER_MM2_PER_M = 8.96; // g/m per mm²
const ALUMINUM_WEIGHT_PER_MM2_PER_M = 2.70; // g/m per mm²

// Materials
const materials = {
    copper: {
        nameKey: "materialCopper",
        resistivity20C: COPPER_RESISTIVITY_20C,
        tempCoefficient: COPPER_TEMP_COEFFICIENT,
        weightPerMm2PerM: COPPER_WEIGHT_PER_MM2_PER_M
    },
    aluminum: {
        nameKey: "materialAluminum",
        resistivity20C: ALUMINUM_RESISTIVITY_20C,
        tempCoefficient: ALUMINUM_TEMP_COEFFICIENT,
        weightPerMm2PerM: ALUMINUM_WEIGHT_PER_MM2_PER_M
    }
};

// Installation method adjustments (°C)
const installationAdjustments = {
    air: 0.0,
    conduit: 10.0,
    isolated: 20.0
};

function getInstallationName(installation) {
    return t(`installation${installation.charAt(0).toUpperCase() + installation.slice(1)}`);
}

// Wire types
const wireTypes = {
    flry: { name: "FLRY", maxTemp: 105.0, description: "Automotive thin-wall PVC" },
    "flry-a": { name: "FLRY-A", maxTemp: 105.0, description: "Automotive flexible stranded" },
    "flry-b": { name: "FLRY-B", maxTemp: 105.0, description: "Automotive symmetrical stranded" },
    thhn: { name: "THHN", maxTemp: 90.0, description: "Thermoplastic, high heat, nylon" },
    thwn: { name: "THWN", maxTemp: 75.0, description: "Thermoplastic, heat/water resistant" },
    xlpe: { name: "XLPE", maxTemp: 90.0, description: "Cross-linked polyethylene" },
    pvc: { name: "PVC", maxTemp: 70.0, description: "Standard PVC" },
    silicon: { name: "Silicone", maxTemp: 200.0, description: "Silicone rubber" },
    generic: { name: "Generic", maxTemp: 90.0, description: "Generic wire" }
};

// Standard metric sizes (mm²)
const standardMetricSizes = [0.5, 0.75, 1.0, 1.5, 2.5, 4.0, 6.0, 10.0, 16.0, 25.0, 35.0, 50.0, 70.0, 95.0, 120.0, 150.0, 185.0, 240.0];

// AWG sizes
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

// Temperature conversion
function fahrenheitToCelsius(f) {
    return (f - 32) * 5 / 9;
}

function celsiusToFahrenheit(c) {
    return c * 9 / 5 + 32;
}

// Calculate resistivity at given temperature
function calculateResistivityAtTemp(material, tempCelsius) {
    return material.resistivity20C * (1 + material.tempCoefficient * (tempCelsius - REFERENCE_TEMP));
}

// Calculate effective operating temperature
function calculateEffectiveTemp(ambientTempCelsius, installation) {
    return ambientTempCelsius + installationAdjustments[installation];
}

// Validate wire temperature
function validateWireTemperature(effectiveTempCelsius, wireType) {
    if (effectiveTempCelsius > wireType.maxTemp) {
        return {
            isValid: false,
            message: t('warningTempExceeds', {
                temp: effectiveTempCelsius.toFixed(1),
                wireType: wireType.name,
                maxTemp: wireType.maxTemp
            })
        };
    }
    if (effectiveTempCelsius > wireType.maxTemp * 0.9) {
        return {
            isValid: true,
            message: t('cautionTempClose', {
                temp: effectiveTempCelsius.toFixed(1),
                wireType: wireType.name,
                maxTemp: wireType.maxTemp
            })
        };
    }
    return { isValid: true, message: "" };
}

// Calculate required cable area
function calculateCableArea(voltage, current, length, maxVoltageDropPercent, material, roundTrip, ambientTempCelsius, installation) {
    const maxVoltageDrop = voltage * (maxVoltageDropPercent / 100.0);
    const distanceFactor = roundTrip ? 2.0 : 1.0;
    const effectiveTemp = calculateEffectiveTemp(ambientTempCelsius, installation);
    const resistivity = calculateResistivityAtTemp(material, effectiveTemp);
    return (current * resistivity * length * distanceFactor) / maxVoltageDrop;
}

// Calculate diameter from area
function areaToDiameter(area) {
    return 2 * Math.sqrt(area / Math.PI);
}

// Calculate cable weight
// weight = area (mm²) × length (m) × weightPerMm2PerM (g/m per mm²)
// For round trip, multiply by 2 (both conductors)
function calculateCableWeight(area, length, material, roundTrip) {
    const weightPerMeter = area * material.weightPerMm2PerM; // g/m
    const totalWeight = weightPerMeter * length; // g
    return roundTrip ? totalWeight * 2 : totalWeight; // g
}

// Find smallest metric size that is >= required area (round up)
function findClosestMetricSize(requiredArea) {
    // Find the smallest standard size that is >= required area
    for (const size of standardMetricSizes) {
        if (size >= requiredArea) {
            return { size: size, diff: size - requiredArea };
        }
    }
    // If required area exceeds all standard sizes, return the largest
    const largestSize = standardMetricSizes[standardMetricSizes.length - 1];
    return { size: largestSize, diff: requiredArea - largestSize };
}

// Find smallest AWG size (largest area) that is >= required area (round up)
function findClosestAWG(requiredArea) {
    // AWG sizes are ordered from smallest to largest area
    // Find the first AWG size where area >= required area
    for (const awg of awgSizes) {
        if (awg.area >= requiredArea) {
            return { label: awg.label, area: awg.area, diff: awg.area - requiredArea };
        }
    }
    // If required area exceeds all AWG sizes, return the largest
    const largestAWG = awgSizes[awgSizes.length - 1];
    return { label: largestAWG.label, area: largestAWG.area, diff: requiredArea - largestAWG.area };
}

// Main calculation function
function calculate() {
    // Get input values
    const voltage = parseFloat(document.getElementById('voltage').value);
    const current = parseFloat(document.getElementById('current').value);
    const length = parseFloat(document.getElementById('length').value);
    const maxVoltageDropPercent = parseFloat(document.getElementById('voltageDrop').value) || 3.0;
    const ambientTemp = parseFloat(document.getElementById('temperature').value) || 20.0;
    const tempUnit = document.getElementById('tempUnit').value;
    const roundTrip = document.getElementById('roundTrip').value === 'true';
    const materialKey = document.getElementById('material').value;
    const installation = document.getElementById('installation').value;
    const wireTypeKey = document.getElementById('wireType').value;

    // Validation
    if (!voltage || voltage <= 0 || voltage > MAX_VOLTAGE) {
        alert(t('errorInvalidVoltage', { maxVoltage: MAX_VOLTAGE }));
        return;
    }
    if (!current || current <= 0) {
        alert(t('errorInvalidCurrent'));
        return;
    }
    if (!length || length <= 0) {
        alert(t('errorInvalidLength'));
        return;
    }

    // Convert temperature to Celsius
    let ambientTempCelsius = ambientTemp;
    if (tempUnit === 'F') {
        ambientTempCelsius = fahrenheitToCelsius(ambientTemp);
    }

    // Get material and wire type
    const material = materials[materialKey];
    const wireType = wireTypes[wireTypeKey];

    // Calculate effective temperature
    const effectiveTemp = calculateEffectiveTemp(ambientTempCelsius, installation);

    // Calculate resistivity at effective temperature
    const resistivity = calculateResistivityAtTemp(material, effectiveTemp);

    // Validate wire temperature
    const tempValidation = validateWireTemperature(effectiveTemp, wireType);

    // Calculate required area
    const requiredArea = calculateCableArea(voltage, current, length, maxVoltageDropPercent, material, roundTrip, ambientTempCelsius, installation);
    const requiredDiameter = areaToDiameter(requiredArea);

    // Find standard sizes
    const metricResult = findClosestMetricSize(requiredArea);
    const awgResult = findClosestAWG(requiredArea);

    // Calculate actual voltage drops
    const distanceFactor = roundTrip ? 2.0 : 1.0;
    const actualDropMetric = (current * resistivity * length * distanceFactor) / metricResult.size;
    const actualDropPercentMetric = (actualDropMetric / voltage) * 100;
    const actualDropAWG = (current * resistivity * length * distanceFactor) / awgResult.area;
    const actualDropPercentAWG = (actualDropAWG / voltage) * 100;

    // Display results
    document.getElementById('result-voltage').textContent = `${voltage.toFixed(1)} V`;
    document.getElementById('result-current').textContent = `${current.toFixed(2)} A`;
    document.getElementById('result-length').textContent = `${length.toFixed(2)} m (${roundTrip ? t('roundTrip') : t('oneWay')})`;
    document.getElementById('result-drop').textContent = `${maxVoltageDropPercent.toFixed(2)}% (${(voltage * maxVoltageDropPercent / 100).toFixed(2)} V)`;
    document.getElementById('result-material').textContent = t(material.nameKey);

    // NEU: Spezifischen Widerstand anzeigen
    document.getElementById('result-resistivity').textContent = `${resistivity.toFixed(4)} Ω·mm²/m`;

    document.getElementById('result-wiretype').textContent = `${wireType.name} (Max: ${wireType.maxTemp}°C)`;
    document.getElementById('result-installation').textContent = getInstallationName(installation);
    document.getElementById('result-ambient').textContent = `${ambientTemp.toFixed(1)}°${tempUnit} (${ambientTempCelsius.toFixed(1)}°C)`;
    document.getElementById('result-effective-temp').textContent = `${effectiveTemp.toFixed(1)}°C`;

    // Display warning if needed
    const warningDiv = document.getElementById('warning');
    if (!tempValidation.isValid || tempValidation.message) {
        warningDiv.classList.remove('hidden');
        if (!tempValidation.isValid) {
            warningDiv.className = 'mb-4 p-4 rounded-lg bg-red-100 dark:bg-red-900/20 border border-red-300 dark:border-red-800';
            warningDiv.innerHTML = `<div class="text-red-800 dark:text-red-200 font-semibold">⚠️ ${tempValidation.message}</div><div class="text-sm text-red-700 dark:text-red-300 mt-2">${t('warningUnsafe')}</div>`;
        } else {
            warningDiv.className = 'mb-4 p-4 rounded-lg bg-yellow-100 dark:bg-yellow-900/20 border border-yellow-300 dark:border-yellow-800';
            warningDiv.innerHTML = `<div class="text-yellow-800 dark:text-yellow-200 font-semibold">⚠️ ${tempValidation.message}</div>`;
        }
    } else {
        warningDiv.classList.add('hidden');
    }

    // Calculate weights
    const requiredWeight = calculateCableWeight(requiredArea, length, material, roundTrip);
    const metricWeight = calculateCableWeight(metricResult.size, length, material, roundTrip);
    const awgWeight = calculateCableWeight(awgResult.area, length, material, roundTrip);

    // Format weight display (show kg if >= 1000g, otherwise g)
    function formatWeight(weightGrams) {
        if (weightGrams >= 1000) {
            return `${(weightGrams / 1000).toFixed(2)} kg`;
        }
        return `${weightGrams.toFixed(1)} g`;
    }

    // Display required size
    document.getElementById('result-area').textContent = `${requiredArea.toFixed(2)} mm²`;
    document.getElementById('result-diameter').textContent = `${requiredDiameter.toFixed(2)} mm`;
    document.getElementById('result-weight').textContent = formatWeight(requiredWeight);

    // Display recommended sizes
    document.getElementById('result-metric').textContent = `${metricResult.size.toFixed(2)} mm²`;
    document.getElementById('result-metric-diff').textContent = metricResult.diff >= 0
        ? `(${t('roundedUp')} ${metricResult.diff.toFixed(2)} mm²)`
        : `(${t('largestAvailable')} ${Math.abs(metricResult.diff).toFixed(2)} mm²)`;
    document.getElementById('result-metric-weight').textContent = `${t('cableWeight')} ${formatWeight(metricWeight)}`;

    document.getElementById('result-awg').textContent = `AWG ${awgResult.label} (${awgResult.area.toFixed(2)} mm²)`;
    document.getElementById('result-awg-diff').textContent = awgResult.diff >= 0
        ? `(${t('roundedUp')} ${awgResult.diff.toFixed(2)} mm²)`
        : `(${t('largestAvailable')} ${Math.abs(awgResult.diff).toFixed(2)} mm²)`;
    document.getElementById('result-awg-weight').textContent = `${t('cableWeight')} ${formatWeight(awgWeight)}`;

    // Display voltage drops
    document.getElementById('result-drop-metric-label').textContent = `${t('with')} ${metricResult.size.toFixed(2)} mm²:`;
    document.getElementById('result-drop-metric').textContent = `${actualDropMetric.toFixed(2)} V (${actualDropPercentMetric.toFixed(2)}%)`;
    document.getElementById('result-drop-awg-label').textContent = `${t('with')} AWG ${awgResult.label} (${awgResult.area.toFixed(2)} mm²):`;
    document.getElementById('result-drop-awg').textContent = `${actualDropAWG.toFixed(2)} V (${actualDropPercentAWG.toFixed(2)}%)`;

    // Show results section
    const resultsSection = document.getElementById('results');
    resultsSection.classList.remove('hidden');

    // Announce results to screen readers
    const srAnnouncement = document.getElementById('sr-announcement');
    const announcement = `${t('calculationResults')}. ${t('requiredCableSize')}: ${requiredArea.toFixed(2)} ${t('crossSectionalArea')}. ${t('recommendedStandardSizes')}: ${t('metric')} ${metricResult.size.toFixed(2)} mm², ${t('awg')} ${awgResult.label}.`;
    srAnnouncement.textContent = announcement;

    // Focus on results section for keyboard users and scroll into view
    resultsSection.setAttribute('tabindex', '-1');
    setTimeout(() => {
        resultsSection.focus();
        resultsSection.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }, 100);
}
