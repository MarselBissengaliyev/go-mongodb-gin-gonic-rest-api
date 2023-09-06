function helloGitty() {
  alert('Hello world')
}

let gitty = document.getElementById("gitty")
gitty.addEventListener("click", () => {
  helloGitty()
})

// ...