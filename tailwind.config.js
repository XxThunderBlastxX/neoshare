/** @type {import('tailwindcss').Config} */
module.exports = {
	content: ["./cmd/web/**/*.html", "./cmd/web/**/*.templ"],
	plugins: [require("daisyui")],
	daisyui: {
		themes: ["luxury"],
	},
	theme: {
		extend: {},
	},
}
