export type Match = {
    id: number;
    homeTeam: string;
    homeAbbr: string;
    homeImg: string;
    awayTeam: string;
    awayAbbr: string;
    awayImg: string;
    date: string;
    time: string;
}

export type Event = {
    id: number;
    matchId: number;
    team: string;
    player: string;
    type: string;
    minute: number;
}
