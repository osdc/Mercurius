# Mercurius
*A Community Newsletter Initiative developed and maintained by OSDC*

## Setting up locally 🚀

- Install hugo on your system (steps vary with different OS). For convenience, refer this [link](https://gohugo.io/installation/).
- Download and install NodeJs from [here](https://nodejs.org/en/download)
- Fork and clone the project and follow these steps
	- `cd Mercurius && npm i` to install npm based dependencies.
	- `hugo serve` to start the local development server (localhost/port details will be available in console).
	- To send out mail notifications, update email list from `scripts/mail-scripts/emails.json` followed by `cd scripts/mail-scripts && go run main.go` (though the email list and credentials needed to authenticate this operation will be confidential to maintainers only).
	- Email template for the mail to be sent is written in HTML format and is super-easy to edit. You can contribute changes to it from `scripts/mail-scripts/email_body.html`,

## Blog Submissions 📖

- Head to `content/post`
- Write the blog that you wish to publish in a md file and save it this folder.
- Be sure to include the header at the start of you blog's md (check `content/post/Example1.md`for more).
- Create a PR and notify the maintainers if your blog doesn't get published within a week.
