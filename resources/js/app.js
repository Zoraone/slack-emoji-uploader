const configApp = new Vue({
    el: '#config',
    data: {
        space: '',
        spaceHasError: false,
        spaceErrorMessage: '',
        token: '',
        tokenHasError: false,
        tokenErrorMessage: '',
        file: null,
        fileHasError: false,
        fileErrorMessage: '',
    },
    methods: {
        getEmojiList() {
            this.spaceHasError = false;
            this.spaceErrorMessage = '';

            this.tokenHasError = false;
            this.tokenErrorMessage = '';

            this.fileHasError = false;
            this.fileErrorMessage = '';

            if (this.space === '') {
                this.spaceHasError = true;
                this.spaceErrorMessage = 'Value Required!';
            }

            if (this.token === '') {
                this.tokenHasError = true;
                this.tokenErrorMessage = 'Value Required!';
            }

            if (this.file === null) {
                this.fileHasError = true;
                this.fileErrorMessage = 'File Required!';
            }

            if (this.spaceHasError && this.tokenHasError && this.fileHasError) { return; }

            astilectron.sendMessage({
                'type': 'get-emoji-list',
                'body': {
                    'token': this.token,
                    'space': this.space,
                    'filename': this.file.path
                }
            }, (response) => {
                console.log(response)
                emojiList.emojis = response;
            });
        },
        onFileChange(event) {
            this.file = event.target.files[0];
        }
    }
});

document.addEventListener('DOMContentLoaded', function () {
    let elems = document.querySelectorAll('.modal');
    M.Modal.init(elems, '');
});

const emojiList = new Vue({
    el: '#emoji-list',
    data: {
        emojis: [],
        searchText: ''
    },
    computed: {
        filteredList() {
            return this.emojis.filter(emoji => {
                return emoji.name.toLowerCase().includes(this.searchText.toLowerCase());
            });
        }
    },
    methods: {
        upload(emoji) {
            emoji.uploadStatus = 'loading';
            this.$forceUpdate();
            astilectron.sendMessage({
                'type': 'upload-emoji',
                'body': {
                    'emoji': emoji
                }
            }, (response) => {
                if (response === 200) {
                    emoji.uploadStatus = 'success';
                } else {
                    emoji.uploadStatus = 'fail';
                }
                this.$forceUpdate();
            });
        },
        uploadAll() {
            for (const emoji of this.emojis) {
                this.upload(emoji);
            }
        }
    }
});