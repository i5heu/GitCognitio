const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
// const MonacoWebpackPlugin = require('monaco-editor-webpack-plugin');
const glob = require('glob');

const htmlFiles = glob.sync('src/**/*.html');
const htmlPlugins = htmlFiles.map(
    (file) => {
        // only use the file name
        bob = file.split('/').pop();

        return new HtmlWebpackPlugin({
            template: file,
            filename: bob,
            inject: bob === 'index.html' ? 'body' : false,
        })
    }
);


module.exports = {
    entry: './src/main.ts',
    mode: 'development',
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                use: 'ts-loader',
                exclude: /node_modules/,
            },
            {
                test: /\.css$/,
                use: ['style-loader', 'css-loader']
            },
            {
                test: /\.s[ac]ss$/i,
                use: [
                    // Creates `style` nodes from JS strings
                    "style-loader",
                    // Translates CSS into CommonJS
                    "css-loader",
                    // Compiles Sass to CSS
                    "sass-loader",
                ],
            },
        ],
    },
    plugins: [
        ...htmlPlugins,
        // new MonacoWebpackPlugin()
    ],
    resolve: {
        extensions: ['.tsx', '.ts', '.js', '...'],
    },
    output: {
        filename: 'bundle.js',
        path: path.resolve(__dirname, 'dist'),
    },
    devServer: {
        static: {
            directory: path.join(__dirname, 'dist'),
        },
        compress: true,
        port: 9000,
    },
};