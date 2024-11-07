export interface QueryProblemListParams {
    page: number;
    page_size: number;
    keyword: string;
    tag: string;
    uid?: number;
}

export interface SubmitForm {
    problem_id: number;
    title: string;
    lang: string;
    code: string;
}

export interface Problem {
    title: string;
    description: string;
    level: number;
    tags: string[];

    id: number;
    status?: number;
    create_at?: string;
    create_by?: number;
    config?: string;
}

interface ProblemConfig {
    testCases: TestCase[];
    compileLimit: Limit;
    runLimit: Limit;
}

interface TestCase {
    input: string;
    output: string;
}

interface Limit {
    cpuLimit: number;
    clockLimit: number;
    memoryLimit: number;
    stackLimit: number;
    procLimit: number;
}

