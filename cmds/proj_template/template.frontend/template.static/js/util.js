function toggleHidden(selector) {
    var elems = document.querySelectorAll(selector);
    elems.forEach(elem => {
        if (elem.classList.contains('hidden')) {
            elem.classList.remove('hidden');
        } else {
            elem.classList.add('hidden');
        }
    });
}

function toggleBetween(selector, subsetClasses1Str, subsetClasses2Str) {
    var elems = document.querySelectorAll(selector);
    const subsetClasses1 = subsetClasses1Str.split(' ');
    const subsetClasses2 = subsetClasses2Str.split(' ');
    elems.forEach(elem => {
        if (subsetClasses1.every(c => elem.classList.contains(c))) {
            subsetClasses1.forEach(c => elem.classList.remove(c));
            subsetClasses2.forEach(c => elem.classList.add(c));
        } else {
            subsetClasses2.forEach(c => elem.classList.remove(c));
            subsetClasses1.forEach(c => elem.classList.add(c));
        }
    });
}

function exclusiveSelectionToggle(targetElem, selector, defaultClassesStr, selectedClassesStr) {
    var elems = document.querySelectorAll(selector);
    const defaultClasses = defaultClassesStr.split(' ');
    const selectedClasses = selectedClassesStr.split(' ');
    elems.forEach(elem => {
        console.log('elem: ', elem, 'targetElem: ', targetElem, 'elem === targetElem: ', elem === targetElem);
        if (elem === targetElem) {
            selectedClasses.forEach(c => elem.classList.add(c));
            defaultClasses.forEach(c => elem.classList.remove(c));
        } else {
            defaultClasses.forEach(c => elem.classList.add(c));
            selectedClasses.forEach(c => elem.classList.remove(c));
        }
    });
}