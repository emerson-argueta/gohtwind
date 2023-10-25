function toggleHidden(elementId) {
    var elem = document.getElementById(elementId);
    if (elem.classList.contains('hidden')) {
        elem.classList.remove('hidden');
    } else {
        elem.classList.add('hidden');
    }
}

function toggleClassesBetween(elementId, subsetClasses1Str, subsetClasses2Str) {
    const subsetClasses1 = subsetClasses1Str.split(' ');
    const subsetClasses2 = subsetClasses2Str.split(' ');
    var elem = document.getElementById(elementId);
    if (subsetClasses1.every(c => elem.classList.contains(c))) {
        subsetClasses1.forEach(c => elem.classList.remove(c));
        subsetClasses2.forEach(c => elem.classList.add(c));
    } else {
        subsetClasses2.forEach(c => elem.classList.remove(c));
        subsetClasses1.forEach(c => elem.classList.add(c));
    }
}
