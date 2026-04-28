export namespace main {
	
	export class AIConfig {
	    provider: string;
	    base_url: string;
	    model: string;
	    api_key: string;
	
	    static createFrom(source: any = {}) {
	        return new AIConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.provider = source["provider"];
	        this.base_url = source["base_url"];
	        this.model = source["model"];
	        this.api_key = source["api_key"];
	    }
	}
	export class StyleCatalogItem {
	    id: string;
	    styleKey: string;
	    label: string;
	    description: string;
	    category: string;
	    source: string;
	
	    static createFrom(source: any = {}) {
	        return new StyleCatalogItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.styleKey = source["styleKey"];
	        this.label = source["label"];
	        this.description = source["description"];
	        this.category = source["category"];
	        this.source = source["source"];
	    }
	}
	export class StyleSkill {
	    style: string;
	    styleName: string;
	    aliases: string[];
	    description: string;
	    category: string;
	    source: string;
	    styleSkill: string;
	
	    static createFrom(source: any = {}) {
	        return new StyleSkill(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.style = source["style"];
	        this.styleName = source["styleName"];
	        this.aliases = source["aliases"];
	        this.description = source["description"];
	        this.category = source["category"];
	        this.source = source["source"];
	        this.styleSkill = source["styleSkill"];
	    }
	}

}

