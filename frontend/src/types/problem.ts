export interface QueryProblemListParams {
    page: number;
    page_size: number;
    keyword: string;
    tag: string;
}

export interface Problem {
    id: number;
    title: string;
    description: string;
    level: number;
    tags: string[];
    status?: number;
}

export interface SubmitForm{
    problem_id: number;
    title: string;
    lang: string;
    code: string;
}