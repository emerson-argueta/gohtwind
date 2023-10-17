function toggleHidden(elementId) {
    var elem = document.getElementById(elementId);
    if (elem.classList.contains('hidden')) {
        elem.classList.remove('hidden');
    } else {
        elem.classList.add('hidden');
    }
}