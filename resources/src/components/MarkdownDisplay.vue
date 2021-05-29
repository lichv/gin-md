
<template>
	<div class="markdown-display">
		<div class="markdown-header">
			
		</div>
		<div class="markdown-content">
			<div class="markdown-side">
				<el-input placeholder="输入关键字进行过滤" v-model="filterText" clearabl></el-input>
				<el-tree ref="tree" :data="items['children']" :props="defaultProps" @node-click="handleNodeClick" :filter-node-method="filterNode"></el-tree>
			</div>
			<div class="markdown-main" id="write">
				<div v-html="html" class="markdown-body"></div>
			</div>
		</div>
	</div>
</template>
<script>
	import {watch} from 'vue'
	import api from '/@/api';
	import marked from 'marked'
	export default {
		name: 'ElMarkdownDisplay',
		data() {
			return {
				filterText: '',
				defaultProps: {
					children: 'children',
					label: 'name'
				},
				drawer: false,
				direction: 'rtl',
				client:null,
				imageUrl: '',
				items:[],
				html:'',
			}
		},
		props: {
			stsUrl: {
				type: String,
				default: "/api/aliyun/sts"
			},
			region: {
				type: String,
				default: "oss-cn-shanghai"
			},
			bucket: {
				type: String,
				default: "easytc"
			},
			rootname: {
				type: String,
				default: "aegicare"
			},
			dirname: {
				type: String,
				default: "picture"
			},
		},
		watch:{
			$router:{
				handler(val,oldval){
					console.log('val',val)
					console.log('oldval',oldval)
				},
			},
			filterText(val) {
				this.$refs.tree.filter(val);
			},
		},

		mounted() {
			let _this = this;
			this.getFiles()
			this.readFile()
			this.$watch(
				() => this.$route.params,
				(toParams, previousParams) => {
					_this.readFile()
				}
				)
		},
		methods:{
			getFiles:function(){
				let _this = this;
				api.getMarkdownFiles.send()
				.then(result => {
					if (result.state==2000) {
						_this.items = result.data
					}
				})
				.catch(e => {
					console.log(e)
				})
			},
			readFile:function(){
				let _this = this;
				var pathMatch = this.$route.params.pathMatch
				if (typeof(pathMatch)!="undefined") {
					var path = pathMatch
					if (pathMatch instanceof Array) {
						var path = pathMatch.join('/')
					}
					api.readMarkdown.send({"filepath":path})
					.then(result => {
						if (result.state==2000) {
							_this.html = marked(result.data)
						}
					})
					.catch(e => {
						console.log(e)
					})
				}else{
					console.log('没进入')
					_this.html = '<div class="welcome"><p>欢迎使用markdown网页浏览工具</p></div>'
				}
			},
			handleNodeClick(data) {
				var relative = data['relative']
				relative = relative.replace(/\\/g,'/')
				if (relative.indexOf('\\')==0 || relative.indexOf('/')==0 ) {
					relative = relative.substr(1)
				}
				if (relative.indexOf('.md') > -1) {
					// this.$router.push({name:'readMarkdown',params:{pathMatch:relative}})
					this.$router.push({path: '/item/'+relative})
				}				
			},
			filterNode(value, data) {
				if (!value) return true;
				return data.path.indexOf(value) !== -1;
			},
		},
	}
</script>
<style>
.markdown-display{
	height: 100%;
	width: 100%;
	padding: 0;
	margin: 0;
}
.markdown-header{
	height: 50px;
	border-bottom: 1px solid #d6d4d4;
}
.markdown-header-handle-btn{
	margin-right: 16px;
	float: right;
	margin-top: 8px;
}
.markdown-content{
	height: calc(100% - 51px);
	display: flex;
}
.markdown-side{
	max-width: 320px;
	height: 100%;
	overflow-y: scroll;
	background-color: #fff;
}
.markdown-main{
	padding: 20px;
	height: 100%;
	overflow-y: scroll;
	flex: 1;
}
</style>