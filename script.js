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

// Get Wikipedia language code
function getWikipediaLang() {
    const langMap = {
        'de': 'de',
        'en': 'en',
        'fr': 'fr',
        'sv': 'sv',
        'de-simple': 'de'
    };
    return langMap[currentLanguage] || 'en';
}

// Update Wikipedia links based on current language
function updateWikipediaLinks() {
    const lang = getWikipediaLang();
    const baseUrl = `https://${lang}.wikipedia.org/wiki/`;

    // Wikipedia article titles in different languages
    const articles = {
        'temp': {
            'de': 'Temperatur',
            'en': 'Temperature',
            'fr': 'Température',
            'sv': 'Temperatur'
        },
        'resistivity': {
            'de': 'Spezifischer_Widerstand',
            'en': 'Electrical_resistivity_and_conductivity',
            'fr': 'Résistivité_électrique',
            'sv': 'Resistivitet'
        },
        'cable-installation': {
            'de': 'Elektrische_Installation',
            'en': 'Electrical_wiring',
            'fr': 'Installation_électrique',
            'sv': 'Elektrisk_installation'
        },
        'voltage-drop': {
            'de': 'Spannungsfall',
            'en': 'Voltage_drop',
            'fr': 'Chute_de_tension',
            'sv': 'Spänningsfall'
        },
        'circle': {
            'de': 'Kreis',
            'en': 'Circle',
            'fr': 'Cercle',
            'sv': 'Cirkel'
        },
        'wire-gauge': {
            'de': 'Leitungsquerschnitt',
            'en': 'American_wire_gauge',
            'fr': 'Calibre_des_fils_électriques',
            'sv': 'Amerikansk_trådmått'
        }
    };

    // Update each link
    const linkMap = {
        'link-wikipedia-temp': articles.temp[lang] || articles.temp.en,
        'link-wikipedia-resistivity': articles.resistivity[lang] || articles.resistivity.en,
        'link-wikipedia-cable-installation': articles['cable-installation'][lang] || articles['cable-installation'].en,
        'link-wikipedia-voltage-drop': articles['voltage-drop'][lang] || articles['voltage-drop'].en,
        'link-wikipedia-circle': articles.circle[lang] || articles.circle.en,
        'link-wikipedia-wire-gauge': articles['wire-gauge'][lang] || articles['wire-gauge'].en,
        'link-wikipedia-voltage-drop-formula': articles['voltage-drop'][lang] || articles['voltage-drop'].en,
        'link-wikipedia-wire-gauge-fuse': articles['wire-gauge'][lang] || articles['wire-gauge'].en
    };

    Object.keys(linkMap).forEach(id => {
        const link = document.getElementById(id);
        if (link) {
            link.href = baseUrl + linkMap[id];
        }
    });
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

    // Re-render MathJax if formulas are visible
    const formulasContent = document.getElementById('formulas-content');
    if (formulasContent && !formulasContent.classList.contains('hidden') && window.MathJax && window.MathJax.typesetPromise) {
        MathJax.typesetPromise([formulasContent]).catch(function (err) {
            console.error('MathJax rendering error:', err);
        });
    }
    
    // Update Wikipedia links
    updateWikipediaLinks();
    
    // Update interactive formula
    if (typeof updateInteractiveFormula === 'function') {
        updateInteractiveFormula();
    }
}

// Constants
const MAX_VOLTAGE = 60.0; // Maximum system voltage (V)
const COPPER_RESISTIVITY_20C = 0.0175; // Ω·mm²/m
const ALUMINUM_RESISTIVITY_20C = 0.0283; // Ω·mm²/m
const COPPER_TEMP_COEFFICIENT = 0.00393; // per °C
const ALUMINUM_TEMP_COEFFICIENT = 0.00403; // per °C
const REFERENCE_TEMP = 25.0; // °C

// Update interactive formula display
// Initialize language on page load
document.addEventListener('DOMContentLoaded', function () {
    // Set language selector
    document.getElementById('language').value = currentLanguage;
    document.documentElement.lang = currentLanguage;

    // Set max voltage attribute
    document.getElementById('voltage').setAttribute('max', MAX_VOLTAGE);

    // Update translations
    updateTranslations();
    
    // Update Wikipedia links
    updateWikipediaLinks();

    // Allow Enter key to trigger calculation and update formula on input
    const inputs = document.querySelectorAll('input[type="number"], select');
    inputs.forEach(input => {
        input.addEventListener('keypress', function (e) {
            if (e.key === 'Enter') {
                calculate();
            }
        });
        
        // Update formula on input change (for number inputs)
        if (input.type === 'number') {
            input.addEventListener('input', function() {
                updateInteractiveFormula();
            });
        }
        
        // Update formula on change (for selects and when input loses focus)
        input.addEventListener('change', function() {
            updateInteractiveFormula();
        });
    });
    
    // Initial formula update (with a small delay to ensure MathJax is loaded)
    setTimeout(function() {
        updateInteractiveFormula();
    }, 100);
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

// Update interactive formula display
function updateInteractiveFormula() {
    const formulaDisplay = document.getElementById('formula-display');
    const formulaValues = document.getElementById('formula-values');
    
    if (!formulaDisplay || !formulaValues) return;
    
    // Get input values
    const voltage = parseFloat(document.getElementById('voltage').value) || 0;
    const current = parseFloat(document.getElementById('current').value) || 0;
    const length = parseFloat(document.getElementById('length').value) || 0;
    const maxVoltageDropPercent = parseFloat(document.getElementById('voltageDrop').value) || 2.0;
    const ambientTemp = parseFloat(document.getElementById('temperature').value) || 20.0;
    const tempUnit = document.getElementById('tempUnit').value;
    const roundTrip = document.getElementById('roundTrip').value === 'true';
    const materialKey = document.getElementById('material').value;
    const installation = document.getElementById('installation').value;
    
    // Convert temperature to Celsius
    let ambientTempCelsius = ambientTemp;
    if (tempUnit === 'F') {
        ambientTempCelsius = fahrenheitToCelsius(ambientTemp);
    }
    
    // Get material
    const material = materials[materialKey];
    if (!material) return;
    
    // Calculate effective temperature
    const effectiveTemp = calculateEffectiveTemp(ambientTempCelsius, installation);
    
    // Calculate resistivity at effective temperature
    const resistivity = calculateResistivityAtTemp(material, effectiveTemp);
    
    // Calculate max voltage drop
    const maxVoltageDrop = voltage > 0 ? voltage * (maxVoltageDropPercent / 100.0) : 0;
    
    // Distance factor
    const distanceFactor = roundTrip ? 2.0 : 1.0;
    
    // Check if we have enough values to show the formula
    const hasValues = voltage > 0 && current > 0 && length > 0 && maxVoltageDrop > 0;
    
    if (hasValues) {
        // Calculate the result
        const area = (current * resistivity * length * distanceFactor) / maxVoltageDrop;
        
        // Build formula with values
        const formula = `A = \\frac{${current.toFixed(2)} \\times ${resistivity.toFixed(4)} \\times ${length.toFixed(2)} \\times ${distanceFactor}}{${maxVoltageDrop.toFixed(2)}} = ${area.toFixed(2)} \\text{ mm}^2`;
        
        formulaDisplay.innerHTML = `\\[${formula}\\]`;
        
        // Show detailed values
        const materialName = t(material.nameKey);
        const distanceText = roundTrip ? t('roundTrip') : t('oneWay');
        const distanceExplanation = roundTrip ? `(${t('roundTrip')})` : `(${t('oneWay')})`;
        
        formulaValues.innerHTML = `
            <div><strong>I</strong> = ${current.toFixed(2)} A</div>
            <div><strong>ρ(T<sub>eff</sub>)</strong> = ${resistivity.toFixed(4)} Ω·mm²/m (${materialName}, ${effectiveTemp.toFixed(1)}°C)</div>
            <div><strong>L</strong> = ${length.toFixed(2)} m ${distanceExplanation}</div>
            <div><strong>d</strong> = ${distanceFactor} ${distanceExplanation}</div>
            <div><strong>ΔV<sub>max</sub></strong> = ${maxVoltageDrop.toFixed(2)} V (${maxVoltageDropPercent.toFixed(1)}% ${t('of')} ${voltage.toFixed(1)} V)</div>
            <div class="mt-2 pt-2 border-t border-gray-300 dark:border-gray-600"><strong>A</strong> = ${area.toFixed(2)} mm²</div>
        `;
        
        // Re-render MathJax
        if (window.MathJax && window.MathJax.typesetPromise) {
            MathJax.typesetPromise([formulaDisplay]).catch(function (err) {
                console.error('MathJax rendering error:', err);
            });
        }
    } else {
        // Show empty formula
        formulaDisplay.innerHTML = `\\[A = \\frac{I \\times \\rho(T_{eff}) \\times L \\times d}{\\Delta V_{max}}\\]`;
        formulaValues.innerHTML = `<div data-i18n="interactiveFormulaNote">Enter values above to see the calculation</div>`;
        
        // Re-render MathJax
        if (window.MathJax && window.MathJax.typesetPromise) {
            MathJax.typesetPromise([formulaDisplay]).catch(function (err) {
                console.error('MathJax rendering error:', err);
            });
        }
    }
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

// Maximum current carrying capacity (Ampacity) for copper wires in automotive applications
// Based on ISO 6722 (Road vehicles - 60 V and 600 V single-core cables), DIN 72551 (Automotive electrical systems),
// and SAE J1127/J1128 (Low voltage primary cable) standards
// Values are for single conductor in free air at 20°C, derated for higher temperatures
// References:
// - ISO 6722: https://www.iso.org/standard/6722.html
// - DIN 72551: https://www.din.de/
// - SAE J1127: https://www.sae.org/standards/content/j1127_201508/
// - SAE J1128: https://www.sae.org/standards/content/j1128_201508/
const ampacityTable = {
    // mm²: base ampacity at 20°C in free air
    // Values based on ISO 6722 and DIN 72551 tables for automotive applications
    0.5: 11,
    0.75: 14,
    1.0: 17,
    1.5: 22,
    2.5: 30,
    4.0: 40,
    6.0: 50,
    10.0: 70,
    16.0: 90,
    25.0: 115,
    35.0: 140,
    50.0: 175,
    70.0: 215,
    95.0: 260,
    120.0: 300,
    150.0: 340,
    185.0: 385,
    240.0: 450
};

// Standard automotive fuse sizes (A)
// Based on ISO 8820 (Road vehicles - Fuse-links) and common automotive blade fuse standards
// References:
// - ISO 8820: https://www.iso.org/standard/8820.html
// - Common sizes: Mini (ATO/ATC), Standard (ATC), Maxi, and ANL/MIDI fuses
const standardFuseSizes = [5, 7.5, 10, 15, 20, 25, 30, 35, 40, 50, 60, 70, 80, 100, 125, 150, 175, 200, 250, 300];

// Temperature derating factors for different ambient temperatures
// Base is 20°C, factors are applied to reduce ampacity at higher temperatures
// Based on IEC 60287 and automotive standards (ISO 6722, DIN 72551)
// References:
// - IEC 60287: Calculation of the continuous current rating of cables
// - ISO 6722: Temperature derating factors for automotive cables
const temperatureDeratingFactors = {
    // temp°C: derating factor
    20: 1.0,
    30: 0.94,
    40: 0.88,
    50: 0.82,
    60: 0.75,
    70: 0.67,
    80: 0.58,
    90: 0.49,
    100: 0.41,
    105: 0.36
};

// Installation derating factors
// Based on cooling conditions: free air provides best cooling, isolated/insulated worst
// References: ISO 6722, DIN 72551, IEC 60287
const installationDeratingFactors = {
    air: 1.0,      // Free air - best cooling
    conduit: 0.8,  // In conduit - reduced cooling
    isolated: 0.7  // Isolated/insulated - worst cooling
};

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

// Calculate maximum current carrying capacity (ampacity) for a given cable
function calculateAmpacity(area, material, installation, ambientTempCelsius, wireType) {
    // Get base ampacity from table (interpolate if needed)
    let baseAmpacity = 0;
    
    // Find closest size in table
    for (const size of standardMetricSizes) {
        if (size >= area) {
            baseAmpacity = ampacityTable[size];
            break;
        }
    }
    
    // If area exceeds table, extrapolate (rough estimate)
    if (baseAmpacity === 0) {
        const largestSize = standardMetricSizes[standardMetricSizes.length - 1];
        const largestAmpacity = ampacityTable[largestSize];
        // Rough linear extrapolation
        baseAmpacity = largestAmpacity * (area / largestSize);
    }
    
    // Apply material factor (aluminum has ~61% of copper conductivity)
    const materialFactor = (material.nameKey === materials.aluminum.nameKey || material === materials.aluminum) ? 0.61 : 1.0;
    
    // Get temperature derating factor (interpolate if needed)
    let tempFactor = 1.0;
    const temps = Object.keys(temperatureDeratingFactors).map(Number).sort((a, b) => a - b);
    
    if (ambientTempCelsius <= temps[0]) {
        tempFactor = temperatureDeratingFactors[temps[0]];
    } else if (ambientTempCelsius >= temps[temps.length - 1]) {
        tempFactor = temperatureDeratingFactors[temps[temps.length - 1]];
    } else {
        // Interpolate between two temperatures
        for (let i = 0; i < temps.length - 1; i++) {
            if (ambientTempCelsius >= temps[i] && ambientTempCelsius <= temps[i + 1]) {
                const t1 = temps[i];
                const t2 = temps[i + 1];
                const f1 = temperatureDeratingFactors[t1];
                const f2 = temperatureDeratingFactors[t2];
                // Linear interpolation
                tempFactor = f1 + (f2 - f1) * ((ambientTempCelsius - t1) / (t2 - t1));
                break;
            }
        }
    }
    
    // Apply installation derating factor
    const installationFactor = installationDeratingFactors[installation] || 1.0;
    
    // Calculate final ampacity
    const ampacity = baseAmpacity * materialFactor * tempFactor * installationFactor;
    
    // Ensure we don't exceed wire type max temperature capability
    // If effective temp is close to max, further derate
    const effectiveTemp = calculateEffectiveTemp(ambientTempCelsius, installation);
    if (effectiveTemp > wireType.maxTemp * 0.9) {
        const tempWarningFactor = Math.max(0.5, 1.0 - ((effectiveTemp - wireType.maxTemp * 0.9) / (wireType.maxTemp * 0.1)));
        return ampacity * tempWarningFactor;
    }
    
    return ampacity;
}

// Find recommended fuse size based on ampacity
function findRecommendedFuse(ampacity, safetyFactor = 0.85) {
    // Apply safety factor (85% = 0.85 means fuse should be 85% of ampacity or less)
    const maxFuseCurrent = ampacity * safetyFactor;
    
    // Find the largest standard fuse that is <= maxFuseCurrent
    let recommendedFuse = null;
    for (const fuseSize of standardFuseSizes) {
        if (fuseSize <= maxFuseCurrent) {
            recommendedFuse = fuseSize;
        } else {
            break;
        }
    }
    
    // If no fuse is small enough, return the smallest
    if (recommendedFuse === null) {
        recommendedFuse = standardFuseSizes[0];
    }
    
    return {
        fuseSize: recommendedFuse,
        maxAmpacity: ampacity,
        safetyFactor: safetyFactor,
        maxFuseCurrent: maxFuseCurrent
    };
}

// Find minimum safe cable size for a given current
function findMinimumSafeCableSize(requiredCurrent, material, installation, ambientTempCelsius, wireType, safetyMargin = 1.1) {
    // Add safety margin (10% = 1.1 means we need 110% of required current capacity)
    const requiredAmpacity = requiredCurrent * safetyMargin;
    
    // Find smallest cable size that can handle the required current
    for (const size of standardMetricSizes) {
        const ampacity = calculateAmpacity(size, material, installation, ambientTempCelsius, wireType);
        if (ampacity >= requiredAmpacity) {
            return {
                size: size,
                area: size,
                ampacity: ampacity,
                safe: true
            };
        }
    }
    
    // If no standard size is large enough, return the largest available
    const largestSize = standardMetricSizes[standardMetricSizes.length - 1];
    const largestAmpacity = calculateAmpacity(largestSize, material, installation, ambientTempCelsius, wireType);
    return {
        size: largestSize,
        area: largestSize,
        ampacity: largestAmpacity,
        safe: largestAmpacity >= requiredAmpacity
    };
}

// Find minimum safe AWG size for a given current
function findMinimumSafeAWGSize(requiredCurrent, material, installation, ambientTempCelsius, wireType, safetyMargin = 1.1) {
    // Add safety margin
    const requiredAmpacity = requiredCurrent * safetyMargin;
    
    // Find smallest AWG size that can handle the required current
    for (const awg of awgSizes) {
        const ampacity = calculateAmpacity(awg.area, material, installation, ambientTempCelsius, wireType);
        if (ampacity >= requiredAmpacity) {
            return {
                label: awg.label,
                area: awg.area,
                ampacity: ampacity,
                safe: true
            };
        }
    }
    
    // If no AWG size is large enough, return the largest available
    const largestAWG = awgSizes[awgSizes.length - 1];
    const largestAmpacity = calculateAmpacity(largestAWG.area, material, installation, ambientTempCelsius, wireType);
    return {
        label: largestAWG.label,
        area: largestAWG.area,
        ampacity: largestAmpacity,
        safe: largestAmpacity >= requiredAmpacity
    };
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
    
    // Calculate ampacity and fuse recommendation for recommended metric size
    const metricAmpacity = calculateAmpacity(metricResult.size, material, installation, ambientTempCelsius, wireType);
    const metricFuseRecommendation = findRecommendedFuse(metricAmpacity);
    
    // Calculate ampacity and fuse recommendation for recommended AWG size
    const awgAmpacity = calculateAmpacity(awgResult.area, material, installation, ambientTempCelsius, wireType);
    const awgFuseRecommendation = findRecommendedFuse(awgAmpacity);

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
    
    // Display fuse recommendations
    document.getElementById('result-fuse-metric').textContent = `${metricFuseRecommendation.fuseSize} A`;
    document.getElementById('result-ampacity-metric').textContent = `${metricAmpacity.toFixed(1)} A`;
    document.getElementById('result-safety-factor-metric').textContent = `${(metricFuseRecommendation.safetyFactor * 100).toFixed(0)}%`;
    
    document.getElementById('result-fuse-awg').textContent = `${awgFuseRecommendation.fuseSize} A`;
    document.getElementById('result-ampacity-awg').textContent = `${awgAmpacity.toFixed(1)} A`;
    document.getElementById('result-safety-factor-awg').textContent = `${(awgFuseRecommendation.safetyFactor * 100).toFixed(0)}%`;
    
    // Check if current exceeds ampacity and show warning with safe cable recommendation
    const fuseWarningDiv = document.getElementById('fuse-warning');
    const fuseWarningText = document.getElementById('fuse-warning-text');
    
    if (current > metricAmpacity || current > awgAmpacity) {
        fuseWarningDiv.classList.remove('hidden');
        fuseWarningDiv.className = 'mt-3 p-3 bg-red-100 dark:bg-red-900/30 rounded border border-red-300 dark:border-red-700';
        
        // Find safe cable sizes
        const safeMetric = findMinimumSafeCableSize(current, material, installation, ambientTempCelsius, wireType);
        const safeAWG = findMinimumSafeAWGSize(current, material, installation, ambientTempCelsius, wireType);
        
        const exceededAmpacity = Math.min(metricAmpacity, awgAmpacity);
        const safeMetricAmpacity = calculateAmpacity(safeMetric.size, material, installation, ambientTempCelsius, wireType);
        const safeAWGAmpacity = calculateAmpacity(safeAWG.area, material, installation, ambientTempCelsius, wireType);
        
        fuseWarningText.innerHTML = `
            <div class="font-semibold mb-2">${t('fuseWarningCurrentExceeds', {
                current: current.toFixed(1),
                ampacity: exceededAmpacity.toFixed(1)
            })}</div>
            <div class="mt-2 text-sm">
                <div class="font-semibold mb-1">${t('recommendedSafeCable')}:</div>
                <div class="ml-2">${t('metric')}: <strong>${safeMetric.size.toFixed(2)} mm²</strong> (${safeMetricAmpacity.toFixed(1)} A ${t('maxCurrentCapacity')})</div>
                <div class="ml-2">${t('awg')}: <strong>AWG ${safeAWG.label}</strong> (${safeAWG.area.toFixed(2)} mm², ${safeAWGAmpacity.toFixed(1)} A ${t('maxCurrentCapacity')})</div>
            </div>
        `;
    } else if (current > metricFuseRecommendation.maxFuseCurrent || current > awgFuseRecommendation.maxFuseCurrent) {
        fuseWarningDiv.classList.remove('hidden');
        fuseWarningDiv.className = 'mt-3 p-3 bg-yellow-100 dark:bg-yellow-900/30 rounded border border-yellow-300 dark:border-yellow-700';
        fuseWarningText.textContent = t('fuseWarningCurrentHigh');
    } else {
        fuseWarningDiv.classList.add('hidden');
    }

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

// Toggle formulas section
function toggleFormulas() {
    const content = document.getElementById('formulas-content');
    const icon = document.getElementById('formulas-icon');
    const toggle = document.getElementById('formulas-toggle');
    const isHidden = content.classList.contains('hidden');

    if (isHidden) {
        content.classList.remove('hidden');
        icon.classList.add('rotate-180');
        toggle.setAttribute('aria-expanded', 'true');
        // Trigger MathJax rendering when formulas are shown
        if (window.MathJax) {
            MathJax.typesetPromise([content]).catch(function (err) {
                console.error('MathJax rendering error:', err);
            });
        }
    } else {
        content.classList.add('hidden');
        icon.classList.remove('rotate-180');
        toggle.setAttribute('aria-expanded', 'false');
    }
}
