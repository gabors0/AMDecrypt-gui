export namespace app {
	
	export class Bento4Settings {
	    mp4decryptPath?: string;
	    binDir?: string;
	
	    static createFrom(source: any = {}) {
	        return new Bento4Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mp4decryptPath = source["mp4decryptPath"];
	        this.binDir = source["binDir"];
	    }
	}
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
	export class Settings {
	    terminal: string;
	    bento4?: Bento4Settings;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.terminal = source["terminal"];
	        this.bento4 = this.convertValues(source["bento4"], Bento4Settings);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

