interface Match {
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

export const getMatches = async () => {
    const matches: Match[] = [
        { id: 1, homeTeam: "Real Madrid", homeAbbr: "RMA", homeImg: "/team_logos/Real Madrid.png", awayTeam: "FC Barcelona", awayAbbr: "BAR", awayImg: "/team_logos/FC Barcelona.png", date: "2023-10-01", time: "18:00" },
        { id: 2, homeTeam: "Atlético de Madrid", homeAbbr: "ATM", homeImg: "/team_logos/Atlético de Madrid.png", awayTeam: "Sevilla FC", awayAbbr: "SEV", awayImg: "/team_logos/Sevilla FC.png", date: "2023-10-02", time: "20:00" },
        { id: 3, homeTeam: "Valencia CF", homeAbbr: "VAL", homeImg: "/team_logos/Valencia CF.png", awayTeam: "Villarreal CF", awayAbbr: "VIL", awayImg: "/team_logos/Villarreal CF.png", date: "2023-10-03", time: "22:00" },
        { id: 4, homeTeam: "Real Sociedad", homeAbbr: "RSO", homeImg: "/team_logos/Real Sociedad.png", awayTeam: "Athletic Bilbao", awayAbbr: "ATH", awayImg: "/team_logos/Athletic Bilbao.png", date: "2023-10-04", time: "18:00" },
        { id: 5, homeTeam: "Real Betis Balompié", homeAbbr: "BET", homeImg: "/team_logos/Real Betis Balompié.png", awayTeam: "Deportivo Alavés", awayAbbr: "ALA", awayImg: "/team_logos/Deportivo Alavés.png", date: "2023-10-05", time: "20:00" },
        { id: 6, homeTeam: "Celta de Vigo", homeAbbr: "CEL", homeImg: "/team_logos/Celta de Vigo.png", awayTeam: "RCD Espanyol Barcelona", awayAbbr: "ESP", awayImg: "/team_logos/RCD Espanyol Barcelona.png", date: "2023-10-06", time: "22:00" }
    ];
    return matches;
}

export const getMatchById = async (id: number): Promise<Match | undefined> => {
    const matches = await getMatches();
    return matches.find(match => match.id === id);
}