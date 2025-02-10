export namespace main {
	
	export class MusicInfo {
	    Size: string;
	    SampleRate: string;
	    Comments: {[key: string]: string};
	    Picture: string;
	
	    static createFrom(source: any = {}) {
	        return new MusicInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Size = source["Size"];
	        this.SampleRate = source["SampleRate"];
	        this.Comments = source["Comments"];
	        this.Picture = source["Picture"];
	    }
	}

}

