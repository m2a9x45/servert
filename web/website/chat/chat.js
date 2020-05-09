const textarea = document.querySelector(".chatBox")

textarea.oninput = function() {
  textarea.style.height = ""; /* Reset the height*/
  textarea.style.height = Math.min(textarea.scrollHeight, 95) + "px";
};