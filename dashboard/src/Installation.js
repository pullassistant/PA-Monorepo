import React, {Component} from 'react';
import './Installation.css';
import Badge from "react-bootstrap/Badge";
import Repository from "./Repository";
import Button from "react-bootstrap/Button";

class Installation extends Component {

    state = {
        repositories: []
    };

    render() {
        let planButtonText;

        if (this.props.config.marketplace_free_plan_id === this.props.installation.plan_id || this.props.config.marketplace_open_source_plan_id === this.props.installation.plan_id) {
            planButtonText = "Upgrade plan"
        } else {
            planButtonText = "Downgrade plan"
        }

        return (
            <div>
                <div className="divider"/>
                <div className="installation">
                    <img className="installation-avatar" src={this.props.installation.avatar_url} alt="avatar"/>
                    <span className="installation-username">{this.props.installation.login}</span>
                    <Badge className="installation-plan" variant="secondary">{this.props.installation.plan_name}</Badge>
                    <div className="installation-actions">
                        <Button href={this.props.installation.html_url} variant="outline-primary" target="_blank">
                            Manage repositories
                        </Button>
                        <Button href={this.props.installation.plan_change_url} variant="outline-primary"
                                target="_blank">
                            {planButtonText}
                        </Button>
                    </div>
                </div>
                {
                    this.state.repositories && this.state.repositories.map((item, key) => <Repository
                        repository={item}/>)
                }
            </div>
        );
    }

    componentDidMount() {
        fetch(process.env.REACT_APP_API_URL + "/repositories/" + this.props.installation.id, {
            method: "GET",
            credentials: 'include'
        })
            .then(res => res.json())
            .then(
                (result) => {
                    this.setState(result)
                },
                (error) => {
                    console.error(error)
                }
            )
    }
}

export default Installation;
