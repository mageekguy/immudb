(window.webpackJsonp=window.webpackJsonp||[]).push([[13],{489:function(e,t,r){"use strict";var n=r(1);var l=function(){var e=arguments.length>0&&void 0!==arguments[0]?arguments[0]:"value",t=arguments.length>1&&void 0!==arguments[1]?arguments[1]:"change";return n.a.extend({name:"proxyable",model:{prop:e,event:t},props:{[e]:{required:!1}},data(){return{internalLazyValue:this[e]}},computed:{internalValue:{get(){return this.internalLazyValue},set(e){e!==this.internalLazyValue&&(this.internalLazyValue=e,this.$emit(t,e))}}},watch:{[e](e){this.internalLazyValue=e}}})}();t.a=l},493:function(e,t,r){},495:function(e,t,r){"use strict";var n=r(1),l=(r(9),r(6),r(5),r(4)),o=(r(493),r(123)),h=r(19),d=r(52),c=r(489),m=r(17),v=r(2),f=r(12);function y(object,e){var t=Object.keys(object);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(object);e&&(r=r.filter((function(e){return Object.getOwnPropertyDescriptor(object,e).enumerable}))),t.push.apply(t,r)}return t}var O=Object(f.a)(h.a,Object(d.b)(["absolute","fixed","top","bottom"]),c.a,m.a).extend({name:"v-progress-linear",props:{active:{type:Boolean,default:!0},backgroundColor:{type:String,default:null},backgroundOpacity:{type:[Number,String],default:null},bufferValue:{type:[Number,String],default:100},color:{type:String,default:"primary"},height:{type:[Number,String],default:4},indeterminate:Boolean,query:Boolean,reverse:Boolean,rounded:Boolean,stream:Boolean,striped:Boolean,value:{type:[Number,String],default:0}},data(){return{internalLazyValue:this.value||0}},computed:{__cachedBackground(){return this.$createElement("div",this.setBackgroundColor(this.backgroundColor||this.color,{staticClass:"v-progress-linear__background",style:this.backgroundStyle}))},__cachedBar(){return this.$createElement(this.computedTransition,[this.__cachedBarType])},__cachedBarType(){return this.indeterminate?this.__cachedIndeterminate:this.__cachedDeterminate},__cachedBuffer(){return this.$createElement("div",{staticClass:"v-progress-linear__buffer",style:this.styles})},__cachedDeterminate(){return this.$createElement("div",this.setBackgroundColor(this.color,{staticClass:"v-progress-linear__determinate",style:{width:Object(v.g)(this.normalizedValue,"%")}}))},__cachedIndeterminate(){return this.$createElement("div",{staticClass:"v-progress-linear__indeterminate",class:{"v-progress-linear__indeterminate--active":this.active}},[this.genProgressBar("long"),this.genProgressBar("short")])},__cachedStream(){return this.stream?this.$createElement("div",this.setTextColor(this.color,{staticClass:"v-progress-linear__stream",style:{width:Object(v.g)(100-this.normalizedBuffer,"%")}})):null},backgroundStyle(){return{opacity:null==this.backgroundOpacity?this.backgroundColor?1:.3:parseFloat(this.backgroundOpacity),[this.isReversed?"right":"left"]:Object(v.g)(this.normalizedValue,"%"),width:Object(v.g)(this.normalizedBuffer-this.normalizedValue,"%")}},classes(){return function(e){for(var i=1;i<arguments.length;i++){var source=null!=arguments[i]?arguments[i]:{};i%2?y(Object(source),!0).forEach((function(t){Object(l.a)(e,t,source[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(source)):y(Object(source)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(source,t))}))}return e}({"v-progress-linear--absolute":this.absolute,"v-progress-linear--fixed":this.fixed,"v-progress-linear--query":this.query,"v-progress-linear--reactive":this.reactive,"v-progress-linear--reverse":this.isReversed,"v-progress-linear--rounded":this.rounded,"v-progress-linear--striped":this.striped},this.themeClasses)},computedTransition(){return this.indeterminate?o.c:o.d},isReversed(){return this.$vuetify.rtl!==this.reverse},normalizedBuffer(){return this.normalize(this.bufferValue)},normalizedValue(){return this.normalize(this.internalLazyValue)},reactive(){return Boolean(this.$listeners.change)},styles(){var e={};return this.active||(e.height=0),this.indeterminate||100===parseFloat(this.normalizedBuffer)||(e.width=Object(v.g)(this.normalizedBuffer,"%")),e}},methods:{genContent(){var slot=Object(v.p)(this,"default",{value:this.internalLazyValue});return slot?this.$createElement("div",{staticClass:"v-progress-linear__content"},slot):null},genListeners(){var e=this.$listeners;return this.reactive&&(e.click=this.onClick),e},genProgressBar(e){return this.$createElement("div",this.setBackgroundColor(this.color,{staticClass:"v-progress-linear__indeterminate",class:{[e]:!0}}))},onClick(e){if(this.reactive){var{width:t}=this.$el.getBoundingClientRect();this.internalValue=e.offsetX/t*100}},normalize:e=>e<0?0:e>100?100:parseFloat(e)},render(e){return e("div",{staticClass:"v-progress-linear",attrs:{role:"progressbar","aria-valuemin":0,"aria-valuemax":this.normalizedBuffer,"aria-valuenow":this.indeterminate?void 0:this.normalizedValue},class:this.classes,style:{bottom:this.bottom?0:void 0,height:this.active?Object(v.g)(this.height):0,top:this.top?0:void 0},on:this.genListeners()},[this.__cachedStream,this.__cachedBackground,this.__cachedBuffer,this.__cachedBar,this.genContent()])}});t.a=n.a.extend().extend({name:"loadable",props:{loading:{type:[Boolean,String],default:!1},loaderHeight:{type:[Number,String],default:2}},methods:{genProgress(){return!1===this.loading?null:this.$slots.progress||this.$createElement(O,{props:{absolute:!0,color:!0===this.loading||""===this.loading?this.color||"primary":this.loading,height:this.loaderHeight,indeterminate:!0}})}}})},500:function(e,t,r){},501:function(e,t,r){"use strict";r.d(t,"a",(function(){return m}));r(9),r(6),r(5);var n=r(4),l=(r(511),r(489)),o=r(17),h=r(12),d=r(11);function c(object,e){var t=Object.keys(object);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(object);e&&(r=r.filter((function(e){return Object.getOwnPropertyDescriptor(object,e).enumerable}))),t.push.apply(t,r)}return t}var m=Object(h.a)(l.a,o.a).extend({name:"base-item-group",props:{activeClass:{type:String,default:"v-item--active"},mandatory:Boolean,max:{type:[Number,String],default:null},multiple:Boolean,tag:{type:String,default:"div"}},data(){return{internalLazyValue:void 0!==this.value?this.value:this.multiple?[]:void 0,items:[]}},computed:{classes(){return function(e){for(var i=1;i<arguments.length;i++){var source=null!=arguments[i]?arguments[i]:{};i%2?c(Object(source),!0).forEach((function(t){Object(n.a)(e,t,source[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(source)):c(Object(source)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(source,t))}))}return e}({"v-item-group":!0},this.themeClasses)},selectedIndex(){return this.selectedItem&&this.items.indexOf(this.selectedItem)||-1},selectedItem(){if(!this.multiple)return this.selectedItems[0]},selectedItems(){return this.items.filter(((e,t)=>this.toggleMethod(this.getValue(e,t))))},selectedValues(){return null==this.internalValue?[]:Array.isArray(this.internalValue)?this.internalValue:[this.internalValue]},toggleMethod(){if(!this.multiple)return e=>this.internalValue===e;var e=this.internalValue;return Array.isArray(e)?t=>e.includes(t):()=>!1}},watch:{internalValue:"updateItemsState",items:"updateItemsState"},created(){this.multiple&&!Array.isArray(this.internalValue)&&Object(d.c)("Model must be bound to an array if the multiple property is true.",this)},methods:{genData(){return{class:this.classes}},getValue:(e,i)=>null==e.value||""===e.value?i:e.value,onClick(e){this.updateInternalValue(this.getValue(e,this.items.indexOf(e)))},register(e){var t=this.items.push(e)-1;e.$on("change",(()=>this.onClick(e))),this.mandatory&&!this.selectedValues.length&&this.updateMandatory(),this.updateItem(e,t)},unregister(e){if(!this._isDestroyed){var t=this.items.indexOf(e),r=this.getValue(e,t);if(this.items.splice(t,1),!(this.selectedValues.indexOf(r)<0)){if(!this.mandatory)return this.updateInternalValue(r);this.multiple&&Array.isArray(this.internalValue)?this.internalValue=this.internalValue.filter((e=>e!==r)):this.internalValue=void 0,this.selectedItems.length||this.updateMandatory(!0)}}},updateItem(e,t){var r=this.getValue(e,t);e.isActive=this.toggleMethod(r)},updateItemsState(){this.$nextTick((()=>{if(this.mandatory&&!this.selectedItems.length)return this.updateMandatory();this.items.forEach(this.updateItem)}))},updateInternalValue(e){this.multiple?this.updateMultiple(e):this.updateSingle(e)},updateMandatory(e){if(this.items.length){var t=this.items.slice();e&&t.reverse();var r=t.find((e=>!e.disabled));if(r){var n=this.items.indexOf(r);this.updateInternalValue(this.getValue(r,n))}}},updateMultiple(e){var t=(Array.isArray(this.internalValue)?this.internalValue:[]).slice(),r=t.findIndex((t=>t===e));this.mandatory&&r>-1&&t.length-1<1||null!=this.max&&r<0&&t.length+1>this.max||(r>-1?t.splice(r,1):t.push(e),this.internalValue=t)},updateSingle(e){var t=e===this.internalValue;this.mandatory&&t||(this.internalValue=t?void 0:e)}},render(e){return e(this.tag,this.genData(),this.$slots.default)}});m.extend({name:"v-item-group",provide(){return{itemGroup:this}}})},510:function(e,t,r){"use strict";r(500)},511:function(e,t,r){},525:function(e,t,r){"use strict";r.r(t);var n=r(20),l={name:"OutputSubNavbar",props:{tab:{type:Number,default:0},tabHasUpdates:{type:Array,default:()=>[]}},data:()=>({mdiText:n.n,mdiAlertCircleOutline:n.a,activeTab:0}),watch:{activeTab(e){this.$emit("update:tab",e)}},methods:{dismissUpdate(data){if(this.tabHasUpdates[data]){var e=this.tabHasUpdates;e[data]=0,this.$emit("update:tabHasUpdates",e)}}}},o=(r(510),r(14)),h=r(18),d=r.n(h),c=r(568),m=r(490),v=r(488),f=r(102),y=r(556),O=r(560),component=Object(o.a)(l,(function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("v-card",{staticClass:"ma-0 pa-0 fill-width bg",attrs:{tile:"",elevation:0}},[r("v-card-text",{staticClass:"ma-0 pa-0",staticStyle:{"box-shadow":"none !important"}},[r("v-tabs",{staticClass:"sub-navbar",attrs:{id:"OutputSubNavbar","slider-color":"primary","show-arrows":"",dense:""},model:{value:e.activeTab,callback:function(t){e.activeTab=t},expression:"activeTab"}},[r("v-tab",{on:{click:function(t){return e.dismissUpdate(0)}}},[r("v-icon",{staticClass:"ml-2 subtitle-1",class:{"gray--text text--darken-1":!e.$vuetify.theme.dark,"gray--text text--lighten-1":e.$vuetify.theme.dark},attrs:{dense:""}},[e._v("\n\t\t\t\t\t"+e._s(e.mdiText)+"\n\t\t\t\t")]),e._v(" "),r("span",{staticClass:"ml-2 body-2 text-capitalize",class:{"gray--text text--darken-1":!e.$vuetify.theme.dark,"gray--text text--lighten-1":e.$vuetify.theme.dark}},[r("v-badge",{attrs:{value:e.tabHasUpdates[0],color:"red lighten-1",content:e.tabHasUpdates[0],bordered:""}},[e._v("\n\t\t\t\t\t\t"+e._s(e.$t("output.code.title"))+"\n\t\t\t\t\t")])],1)],1)],1)],1)],1)}),[],!1,null,null,null);t.default=component.exports;d()(component,{VBadge:c.a,VCard:m.a,VCardText:v.a,VIcon:f.a,VTab:y.a,VTabs:O.a})}}]);