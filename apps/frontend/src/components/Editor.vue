<template>

    <div>
        <div id="editor">
            <textarea rows="50" cols="50" id="markdown" v-model="markdown" @change="updateHTML"></textarea>
        </div>

        <div>
            <p>renders...</p>
            <pre>{{html}}</pre>
        </div>
    </div>
</template>

<script>
    import myMarked from 'marked';

    export default {
        props: {
        },
        data() {
            return {
                marked: function() {
                    myMarked.setOptions({
                        baseUrl: null,
                        breaks: false,
                        gfm: true,
                        headerIds: true,
                        headerPrefix: "",
                        highlight: null,
                        langPrefix: "language-",
                        mangle: true,
                        pedantic: false,
                        sanitize: false,
                        sanitizer: null,
                        silent: false,
                        smartLists: false,
                        smartypants: false,
                        tables: true,
                        xhtml: false
                    });
                },
                markdown: "some test",
                html: "<h1>some html</h1>",
            };
        },
        methods: {
            updateHTML() {
                this.html = myMarked(this.markdown);
            },
        }
    }
</script>

<style>

</style>
