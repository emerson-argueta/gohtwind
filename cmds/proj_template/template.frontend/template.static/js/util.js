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
    const defaultClasses = defaultClassesStr.split(' ').filter(c => c);
    const selectedClasses = selectedClassesStr.split(' ').filter(c => c);
    elems.forEach(elem => {
        if (elem === targetElem) {
            selectedClasses.forEach(c => elem.classList.add(c));
            defaultClasses.forEach(c => elem.classList.remove(c));
        } else {
            defaultClasses.forEach(c => elem.classList.add(c));
            selectedClasses.forEach(c => elem.classList.remove(c));
        }
    });
}

function handleEscapeKey(event, selector) {
    if (event.keyCode === 27) { // 27 is the key code for Escape
        var containers = document.querySelectorAll(selector);
        containers.forEach(container => {
            if (!container.classList.contains('hidden')) {
                container.classList.add('hidden');
            }
        });
    }
}