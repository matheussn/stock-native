export namespace application {
	
	export class BasketSubstitutionInput {
	    originalProductId: string;
	    productId: string;
	    quantityPerBasket: string;
	
	    static createFrom(source: any = {}) {
	        return new BasketSubstitutionInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.originalProductId = source["originalProductId"];
	        this.productId = source["productId"];
	        this.quantityPerBasket = source["quantityPerBasket"];
	    }
	}
	export class BasketTemplateItemInput {
	    productId: string;
	    quantityPerBasket: string;
	
	    static createFrom(source: any = {}) {
	        return new BasketTemplateItemInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.productId = source["productId"];
	        this.quantityPerBasket = source["quantityPerBasket"];
	    }
	}
	export class CreateBasketMovementInput {
	    basketTemplateId: string;
	    basketsCount: string;
	    sourceId: string;
	    // Go type: time
	    movementDate: any;
	    note: string;
	    substitutions: BasketSubstitutionInput[];
	
	    static createFrom(source: any = {}) {
	        return new CreateBasketMovementInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.basketTemplateId = source["basketTemplateId"];
	        this.basketsCount = source["basketsCount"];
	        this.sourceId = source["sourceId"];
	        this.movementDate = this.convertValues(source["movementDate"], null);
	        this.note = source["note"];
	        this.substitutions = this.convertValues(source["substitutions"], BasketSubstitutionInput);
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
	export class CreateBasketTemplateInput {
	    name: string;
	    items: BasketTemplateItemInput[];
	
	    static createFrom(source: any = {}) {
	        return new CreateBasketTemplateInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.items = this.convertValues(source["items"], BasketTemplateItemInput);
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
	export class CreateCategoryInput {
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateCategoryInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	    }
	}
	export class CreateDestinationInput {
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateDestinationInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	    }
	}
	export class CreateFamilyInput {
	    name: string;
	    cestasPerPeriod: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateFamilyInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.cestasPerPeriod = source["cestasPerPeriod"];
	    }
	}
	export class MovementBatchItemInput {
	    productId: string;
	    quantity: string;
	
	    static createFrom(source: any = {}) {
	        return new MovementBatchItemInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.productId = source["productId"];
	        this.quantity = source["quantity"];
	    }
	}
	export class CreateMovementBatchInput {
	    // Go type: time
	    movementDate: any;
	    sourceId: string;
	    note: string;
	    items: MovementBatchItemInput[];
	
	    static createFrom(source: any = {}) {
	        return new CreateMovementBatchInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.movementDate = this.convertValues(source["movementDate"], null);
	        this.sourceId = source["sourceId"];
	        this.note = source["note"];
	        this.items = this.convertValues(source["items"], MovementBatchItemInput);
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
	export class CreateMovementInput {
	    productId: string;
	    quantity: string;
	    // Go type: time
	    movementDate: any;
	    sourceId: string;
	    note: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateMovementInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.productId = source["productId"];
	        this.quantity = source["quantity"];
	        this.movementDate = this.convertValues(source["movementDate"], null);
	        this.sourceId = source["sourceId"];
	        this.note = source["note"];
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
	export class CreateOriginInput {
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateOriginInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	    }
	}
	export class CreateProductInput {
	    name: string;
	    categoryId: string;
	    measureUnit: string;
	    packageAmount: string;
	    lowStockThreshold: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateProductInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.categoryId = source["categoryId"];
	        this.measureUnit = source["measureUnit"];
	        this.packageAmount = source["packageAmount"];
	        this.lowStockThreshold = source["lowStockThreshold"];
	    }
	}
	export class CurrentStockFilter {
	    categoryId: string;
	    search: string;
	
	    static createFrom(source: any = {}) {
	        return new CurrentStockFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.categoryId = source["categoryId"];
	        this.search = source["search"];
	    }
	}
	export class ImportBackupInput {
	    filePath: string;
	    confirmReplace: boolean;
	    confirmUnderstand: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ImportBackupInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filePath = source["filePath"];
	        this.confirmReplace = source["confirmReplace"];
	        this.confirmUnderstand = source["confirmUnderstand"];
	    }
	}
	export class ListMovementsFilter {
	    type: string;
	    productId: string;
	    // Go type: time
	    dateFrom?: any;
	    // Go type: time
	    dateTo?: any;
	
	    static createFrom(source: any = {}) {
	        return new ListMovementsFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.productId = source["productId"];
	        this.dateFrom = this.convertValues(source["dateFrom"], null);
	        this.dateTo = this.convertValues(source["dateTo"], null);
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
	
	export class ReportsFilter {
	    // Go type: time
	    dateFrom: any;
	    // Go type: time
	    dateTo: any;
	
	    static createFrom(source: any = {}) {
	        return new ReportsFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dateFrom = this.convertValues(source["dateFrom"], null);
	        this.dateTo = this.convertValues(source["dateTo"], null);
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
	export class StockPlanningFilter {
	    basketTemplateId: string;
	    horizonDays: number;
	
	    static createFrom(source: any = {}) {
	        return new StockPlanningFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.basketTemplateId = source["basketTemplateId"];
	        this.horizonDays = source["horizonDays"];
	    }
	}
	export class UpdateBasketTemplateInput {
	    id: string;
	    name: string;
	    items: BasketTemplateItemInput[];
	
	    static createFrom(source: any = {}) {
	        return new UpdateBasketTemplateInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.items = this.convertValues(source["items"], BasketTemplateItemInput);
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
	export class UpdateCategoryInput {
	    id: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateCategoryInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}
	export class UpdateDestinationInput {
	    id: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateDestinationInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}
	export class UpdateFamilyInput {
	    id: string;
	    name: string;
	    cestasPerPeriod: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateFamilyInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.cestasPerPeriod = source["cestasPerPeriod"];
	    }
	}
	export class UpdateOriginInput {
	    id: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateOriginInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}
	export class UpdateProductInput {
	    id: string;
	    name: string;
	    categoryId: string;
	    measureUnit: string;
	    packageAmount: string;
	    lowStockThreshold: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateProductInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.categoryId = source["categoryId"];
	        this.measureUnit = source["measureUnit"];
	        this.packageAmount = source["packageAmount"];
	        this.lowStockThreshold = source["lowStockThreshold"];
	    }
	}

}

export namespace domain {
	
	export class BasketTemplateItem {
	    id: string;
	    basketTemplateId: string;
	    productId: string;
	    productName?: string;
	    measureUnit: string;
	    quantityPerBasket: number;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new BasketTemplateItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.basketTemplateId = source["basketTemplateId"];
	        this.productId = source["productId"];
	        this.productName = source["productName"];
	        this.measureUnit = source["measureUnit"];
	        this.quantityPerBasket = source["quantityPerBasket"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
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
	export class BasketTemplate {
	    id: string;
	    name: string;
	    active: boolean;
	    items: BasketTemplateItem[];
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new BasketTemplate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.active = source["active"];
	        this.items = this.convertValues(source["items"], BasketTemplateItem);
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
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
	
	export class Category {
	    id: string;
	    name: string;
	    active: boolean;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Category(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.active = source["active"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
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
	export class Destination {
	    id: string;
	    name: string;
	    active: boolean;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Destination(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.active = source["active"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
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
	export class Family {
	    id: string;
	    name: string;
	    cestasPerPeriod: number;
	    active: boolean;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Family(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.cestasPerPeriod = source["cestasPerPeriod"];
	        this.active = source["active"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
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
	export class GroupedTotal {
	    id: string;
	    name: string;
	    quantity: number;
	
	    static createFrom(source: any = {}) {
	        return new GroupedTotal(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.quantity = source["quantity"];
	    }
	}
	export class Movement {
	    id: string;
	    type: string;
	    productId: string;
	    productName?: string;
	    quantity: number;
	    // Go type: time
	    movementDate: any;
	    originId?: string;
	    originName?: string;
	    destinationId?: string;
	    destinationName?: string;
	    note: string;
	    // Go type: time
	    createdAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Movement(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.productId = source["productId"];
	        this.productName = source["productName"];
	        this.quantity = source["quantity"];
	        this.movementDate = this.convertValues(source["movementDate"], null);
	        this.originId = source["originId"];
	        this.originName = source["originName"];
	        this.destinationId = source["destinationId"];
	        this.destinationName = source["destinationName"];
	        this.note = source["note"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
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
	export class Origin {
	    id: string;
	    name: string;
	    active: boolean;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Origin(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.active = source["active"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
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
	export class Product {
	    id: string;
	    name: string;
	    categoryId: string;
	    measureUnit: string;
	    packageAmount: number;
	    lowStockThreshold: number;
	    active: boolean;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Product(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.categoryId = source["categoryId"];
	        this.measureUnit = source["measureUnit"];
	        this.packageAmount = source["packageAmount"];
	        this.lowStockThreshold = source["lowStockThreshold"];
	        this.active = source["active"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
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
	export class Reports {
	    entriesByPeriod: Movement[];
	    entriesByOrigin: GroupedTotal[];
	    exitsByPeriod: Movement[];
	    exitsByDestination: GroupedTotal[];
	    movementHistory: Movement[];
	
	    static createFrom(source: any = {}) {
	        return new Reports(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entriesByPeriod = this.convertValues(source["entriesByPeriod"], Movement);
	        this.entriesByOrigin = this.convertValues(source["entriesByOrigin"], GroupedTotal);
	        this.exitsByPeriod = this.convertValues(source["exitsByPeriod"], Movement);
	        this.exitsByDestination = this.convertValues(source["exitsByDestination"], GroupedTotal);
	        this.movementHistory = this.convertValues(source["movementHistory"], Movement);
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
	export class StockItem {
	    productId: string;
	    productName: string;
	    categoryId: string;
	    categoryName: string;
	    measureUnit: string;
	    packageAmount: number;
	    lowStockThreshold: number;
	    currentStock: number;
	    isLowStock: boolean;
	
	    static createFrom(source: any = {}) {
	        return new StockItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.productId = source["productId"];
	        this.productName = source["productName"];
	        this.categoryId = source["categoryId"];
	        this.categoryName = source["categoryName"];
	        this.measureUnit = source["measureUnit"];
	        this.packageAmount = source["packageAmount"];
	        this.lowStockThreshold = source["lowStockThreshold"];
	        this.currentStock = source["currentStock"];
	        this.isLowStock = source["isLowStock"];
	    }
	}
	export class StockPlanningItem {
	    productId: string;
	    productName: string;
	    measureUnit: string;
	    currentStock: number;
	    requiredStock: number;
	    projectedBalance: number;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new StockPlanningItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.productId = source["productId"];
	        this.productName = source["productName"];
	        this.measureUnit = source["measureUnit"];
	        this.currentStock = source["currentStock"];
	        this.requiredStock = source["requiredStock"];
	        this.projectedBalance = source["projectedBalance"];
	        this.status = source["status"];
	    }
	}
	export class StockPlanning {
	    basketTemplateId: string;
	    basketName: string;
	    horizonDays: number;
	    requiredBaskets: number;
	    items: StockPlanningItem[];
	
	    static createFrom(source: any = {}) {
	        return new StockPlanning(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.basketTemplateId = source["basketTemplateId"];
	        this.basketName = source["basketName"];
	        this.horizonDays = source["horizonDays"];
	        this.requiredBaskets = source["requiredBaskets"];
	        this.items = this.convertValues(source["items"], StockPlanningItem);
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

