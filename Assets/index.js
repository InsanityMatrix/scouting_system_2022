//Pregame
function teamInputs() {
	setTimeout(function() {
		var num = parseInt(document.getElementById("teamNum").value,10);
		if (!((num >= 0) || (num < 10000))) {
			document.getElementById("teamNum").value = "";
			document.getElementById("pregameBadTime").classList.add("fadeOut");
		}
	}, 10);
	setTimeout(function() {
		document.getElementById("pregameBadTime").classList.remove("fadeOut");
	}, 3010);
}

function matchInputs() {
	setTimeout(function() {
		var num = parseInt(document.getElementById("matchNum").value,10);
		if (!((num >= 0) || (num < 10000))) {
			document.getElementById("matchNum").value = "";
			document.getElementById("pregameBadTime").classList.add("fadeOut");
		}
	}, 10);
	setTimeout(function() {
		document.getElementById("pregameBadTime").classList.remove("fadeOut");
	}, 3010);
}

function checkDone() {
	if (((document.getElementById("preloadYes").checked) || (document.getElementById("preloadNo").checked)) && !(document.getElementById("station").value == "")) {
		if (!(document.getElementById("teamNum").value == "") && !(document.getElementById("matchNum").value == "")) {
			document.getElementById("tabSwitch").style.opacity = 1;
		}
	}
}
//Endgame
function noAttempt() {
	if (document.getElementById("noAttempt").checked == true) {
		document.getElementById("climb").hidden = true;
	} else {
		document.getElementById("climb").hidden = false;
	}
}
function disable(ident) {
	var ident2 = "s" + ident.substring(1, ident.length);
	if (document.getElementById(ident).checked == true) {
		document.getElementById(ident2).disabled = false;
	} else {
		document.getElementById(ident2).disabled = true;
	}
}