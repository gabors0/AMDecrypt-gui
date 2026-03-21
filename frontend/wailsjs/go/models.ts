export namespace app {
	
	export class InstanceConfig {
	    url: string;
	    secure: boolean;
	
	    static createFrom(source: any = {}) {
	        return new InstanceConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.secure = source["secure"];
	    }
	}

}

