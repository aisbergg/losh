// Input mask plugin

import IMask from 'imask';

var maskElementList = [].slice.call(document.querySelectorAll('[data-mask]'));
maskElementList.map(function (maskEl) {
	return new IMask(maskEl, {
		mask: maskEl.dataset.mask,
		lazy: maskEl.dataset['mask-visible'] === 'true'
	})
});

var numMaskElementList = [].slice.call(document.querySelectorAll('[data-number-mask]'));
numMaskElementList.map(function (maskEl) {
	return new IMask(maskEl, {
		mask: Number,
        min: maskEl.dataset['number-mask-min'],
        max: maskEl.dataset['number-mask-max'],
	})
});
