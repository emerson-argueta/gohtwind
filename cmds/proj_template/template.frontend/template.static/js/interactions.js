/**
 * @param {Element|string} targetElem Element or string of selector to match element to be toggled
 * @param {string} selector String of selector to match elements to be toggled
 * @param {string} defaultClassesStr String of classes to be added to all elements matching selector and removed from targetElem
 * @param {string} selectedClassesStr String of classes to be added to targetElem and removed from all other elements matching selector
 */
function exclusiveSelectionToggle(targetElem, selector, defaultClassesStr, selectedClassesStr) {
    if (typeof targetElem === 'string') {
        targetElem = document.querySelector(targetElem);
    }
    const elems = document.querySelectorAll(selector);
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


/**
 * @param {Element|string} targetElem Element or string of selector to match element to be toggled
 * @param {string} classes1Str String of classes to be removed from targetElem
 * @param {string} classes2Str String of classes to be added to targetElem
 */
function classesSwap(targetElem, classes1Str, classes2Str) {
    if (typeof targetElem === 'string') {
        targetElem = document.querySelector(targetElem);
    }
    const classes1 = classes1Str.split(' ').filter(c => c);
    const classes2 = classes2Str.split(' ').filter(c => c);
    classes1.forEach(c => {
        if (targetElem.classList.contains(c)) {
            targetElem.classList.remove(c);
        }
    });
    classes2.forEach(c => {
        if (!targetElem.classList.contains(c)) {
            targetElem.classList.add(c);
        }
    });
}

/**
 * @param {NodeList<Element>|string} targetElems Elements or string of selector to match elements to be toggled
 * @param {string} classStr String of classes to be toggled
 */
function toggleClasses(targetElems, classStr) {
    if (typeof targetElems === 'string') {
        targetElems = document.querySelectorAll(targetElems);
    }
    const classes = classStr.split(' ').filter(c => c);
    targetElems.forEach(elem => {
        classes.forEach(c => {
            elem.classList.toggle(c);
        });
    });
}

/**
 * @param {Element|string} targetElem Element or string of selector to match element to be toggled
 * @param {string} classStr String of classes to be removed
 */
function removeClasses(targetElem, classStr) {
    if (typeof targetElem === 'string') {
        targetElem = document.querySelector(targetElem);
    }
    const classes = classStr.split(' ').filter(c => c);
    classes.forEach(c => {
        targetElem.classList.remove(c);
    });
}

/**
 * @param {Element|string} targetElem Element or string of selector to match element to be toggled
 * @param {string} classStr String of classes to be added
 */
function addClasses(targetElem, classStr) {
    if (typeof targetElem === 'string') {
        targetElem = document.querySelector(targetElem);
    }
    const classes = classStr.split(' ').filter(c => c);
    classes.forEach(c => {
        if (!targetElem.classList.contains(c)) {
            targetElem.classList.add(c);
        }
    });
}

/**
 * @param {Event} event Event object
 * @param {String} selector String of selector to match elements to be hidden
 */
function hideOnEscape(event, selector) {
    if (event.keyCode === 27) {
        var containers = document.querySelectorAll(selector);
        containers.forEach(container => {
            if (!container.classList.contains('hidden')) {
                container.classList.add('hidden');
            }
        });
    }
}