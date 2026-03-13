export namespace db {
	
	export class RestoreResult {
	    restored_path: string;
	    previous_backup_path: string;
	
	    static createFrom(source: any = {}) {
	        return new RestoreResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.restored_path = source["restored_path"];
	        this.previous_backup_path = source["previous_backup_path"];
	    }
	}

}

export namespace main {
	
	export class AppStartupStatus {
	    ready: boolean;
	    message: string;
	    technical_error: string;
	    database_path: string;
	
	    static createFrom(source: any = {}) {
	        return new AppStartupStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ready = source["ready"];
	        this.message = source["message"];
	        this.technical_error = source["technical_error"];
	        this.database_path = source["database_path"];
	    }
	}

}

export namespace models {
	
	export class AssistentialWork {
	    id: number;
	    name: string;
	    description: string;
	    is_active: number;
	    created_at: string;
	
	    static createFrom(source: any = {}) {
	        return new AssistentialWork(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.is_active = source["is_active"];
	        this.created_at = source["created_at"];
	    }
	}
	export class Family {
	    id: number;
	    assistential_work_id: number;
	    name: string;
	    member_count: number;
	    address: string;
	    contact: string;
	    is_active: number;
	    created_at: string;
	
	    static createFrom(source: any = {}) {
	        return new Family(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.assistential_work_id = source["assistential_work_id"];
	        this.name = source["name"];
	        this.member_count = source["member_count"];
	        this.address = source["address"];
	        this.contact = source["contact"];
	        this.is_active = source["is_active"];
	        this.created_at = source["created_at"];
	    }
	}
	export class FamilyGroupAssignment {
	    id: number;
	    family_id: number;
	    product_group_id: number;
	    started_at: string;
	    ended_at: string;
	
	    static createFrom(source: any = {}) {
	        return new FamilyGroupAssignment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.family_id = source["family_id"];
	        this.product_group_id = source["product_group_id"];
	        this.started_at = source["started_at"];
	        this.ended_at = source["ended_at"];
	    }
	}
	export class Movement {
	    id: number;
	    assistential_work_id: number;
	    type: string;
	    notes: string;
	    created_at: string;
	
	    static createFrom(source: any = {}) {
	        return new Movement(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.assistential_work_id = source["assistential_work_id"];
	        this.type = source["type"];
	        this.notes = source["notes"];
	        this.created_at = source["created_at"];
	    }
	}
	export class MovementGroupItem {
	    id: number;
	    movement_id: number;
	    product_group_id: number;
	    quantity: number;
	
	    static createFrom(source: any = {}) {
	        return new MovementGroupItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.movement_id = source["movement_id"];
	        this.product_group_id = source["product_group_id"];
	        this.quantity = source["quantity"];
	    }
	}
	export class MovementGroupItemResolution {
	    id: number;
	    movement_group_item_id: number;
	    product_variation_id: number;
	    quantity: number;
	
	    static createFrom(source: any = {}) {
	        return new MovementGroupItemResolution(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.movement_group_item_id = source["movement_group_item_id"];
	        this.product_variation_id = source["product_variation_id"];
	        this.quantity = source["quantity"];
	    }
	}
	export class MovementProductItem {
	    id: number;
	    movement_id: number;
	    product_variation_id: number;
	    quantity: number;
	
	    static createFrom(source: any = {}) {
	        return new MovementProductItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.movement_id = source["movement_id"];
	        this.product_variation_id = source["product_variation_id"];
	        this.quantity = source["quantity"];
	    }
	}
	export class Product {
	    id: number;
	    name: string;
	    base_unit: string;
	    description: string;
	    is_active: number;
	    created_at: string;
	
	    static createFrom(source: any = {}) {
	        return new Product(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.base_unit = source["base_unit"];
	        this.description = source["description"];
	        this.is_active = source["is_active"];
	        this.created_at = source["created_at"];
	    }
	}
	export class ProductGroup {
	    id: number;
	    name: string;
	    description: string;
	    is_active: number;
	    created_at: string;
	
	    static createFrom(source: any = {}) {
	        return new ProductGroup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.is_active = source["is_active"];
	        this.created_at = source["created_at"];
	    }
	}
	export class ProductGroupItem {
	    id: number;
	    product_group_id: number;
	    product_id: number;
	    base_quantity: number;
	
	    static createFrom(source: any = {}) {
	        return new ProductGroupItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.product_group_id = source["product_group_id"];
	        this.product_id = source["product_id"];
	        this.base_quantity = source["base_quantity"];
	    }
	}
	export class ProductVariation {
	    id: number;
	    product_id: number;
	    description: string;
	    base_quantity: number;
	    current_stock: number;
	    is_active: number;
	    created_at: string;
	
	    static createFrom(source: any = {}) {
	        return new ProductVariation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.product_id = source["product_id"];
	        this.description = source["description"];
	        this.base_quantity = source["base_quantity"];
	        this.current_stock = source["current_stock"];
	        this.is_active = source["is_active"];
	        this.created_at = source["created_at"];
	    }
	}

}

export namespace services {
	
	export class AssistentialWorkOutflowKgItem {
	    assistential_work_id: number;
	    assistential_work_name: string;
	    output_kg: number;
	
	    static createFrom(source: any = {}) {
	        return new AssistentialWorkOutflowKgItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.assistential_work_id = source["assistential_work_id"];
	        this.assistential_work_name = source["assistential_work_name"];
	        this.output_kg = source["output_kg"];
	    }
	}
	export class MovementGroupItemInput {
	    product_group_id: number;
	    quantity: number;
	
	    static createFrom(source: any = {}) {
	        return new MovementGroupItemInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.product_group_id = source["product_group_id"];
	        this.quantity = source["quantity"];
	    }
	}
	export class MovementProductItemInput {
	    product_variation_id: number;
	    quantity: number;
	
	    static createFrom(source: any = {}) {
	        return new MovementProductItemInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.product_variation_id = source["product_variation_id"];
	        this.quantity = source["quantity"];
	    }
	}
	export class CreateMovementInput {
	    assistential_work_id: number;
	    type: string;
	    notes: string;
	    product_items: MovementProductItemInput[];
	    group_items: MovementGroupItemInput[];
	
	    static createFrom(source: any = {}) {
	        return new CreateMovementInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.assistential_work_id = source["assistential_work_id"];
	        this.type = source["type"];
	        this.notes = source["notes"];
	        this.product_items = this.convertValues(source["product_items"], MovementProductItemInput);
	        this.group_items = this.convertValues(source["group_items"], MovementGroupItemInput);
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
	export class MonthlyDemandItem {
	    product_id: number;
	    product_name: string;
	    base_unit: string;
	    required_base_quantity: number;
	
	    static createFrom(source: any = {}) {
	        return new MonthlyDemandItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.product_id = source["product_id"];
	        this.product_name = source["product_name"];
	        this.base_unit = source["base_unit"];
	        this.required_base_quantity = source["required_base_quantity"];
	    }
	}
	export class MonthlyMovementFlowKgItem {
	    month_key: string;
	    input_kg: number;
	    output_kg: number;
	
	    static createFrom(source: any = {}) {
	        return new MonthlyMovementFlowKgItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.month_key = source["month_key"];
	        this.input_kg = source["input_kg"];
	        this.output_kg = source["output_kg"];
	    }
	}
	
	
	export class StockCoverageItem {
	    product_id: number;
	    product_name: string;
	    base_unit: string;
	    required_base_quantity: number;
	    available_base_quantity: number;
	    shortfall_base_quantity: number;
	    is_covered: boolean;
	
	    static createFrom(source: any = {}) {
	        return new StockCoverageItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.product_id = source["product_id"];
	        this.product_name = source["product_name"];
	        this.base_unit = source["base_unit"];
	        this.required_base_quantity = source["required_base_quantity"];
	        this.available_base_quantity = source["available_base_quantity"];
	        this.shortfall_base_quantity = source["shortfall_base_quantity"];
	        this.is_covered = source["is_covered"];
	    }
	}
	export class VariationStockStatus {
	    product_id: number;
	    product_name: string;
	    base_unit: string;
	    variation_id: number;
	    variation_description: string;
	    base_quantity: number;
	    current_stock: number;
	
	    static createFrom(source: any = {}) {
	        return new VariationStockStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.product_id = source["product_id"];
	        this.product_name = source["product_name"];
	        this.base_unit = source["base_unit"];
	        this.variation_id = source["variation_id"];
	        this.variation_description = source["variation_description"];
	        this.base_quantity = source["base_quantity"];
	        this.current_stock = source["current_stock"];
	    }
	}

}

